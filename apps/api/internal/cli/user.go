package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/google/uuid"
	"github.com/orvexcc/billing/api/internal/db"
	"github.com/orvexcc/billing/api/internal/types"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "User management commands",
	Long:  `Commands for managing users in the billing system.`,
}

var makeAdminCmd = &cobra.Command{
	Use:   "make-admin [email_or_username]",
	Short: "Make a user an admin",
	Long:  `Promote a user to admin role by their email address or username.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		identifier := args[0]
		
		var user types.User
		result := db.DB.Where("email = ? OR username = ?", identifier, identifier).First(&user)
		if result.Error != nil {
			fmt.Printf("Error: User with email/username '%s' not found\n", identifier)
			return
		}

		if user.Role == types.UserRoleAdmin || user.Role == types.UserRoleSuper {
			fmt.Printf("User %s (%s) is already an admin\n", user.Username, user.Email)
			return
		}

		user.Role = types.UserRoleAdmin
		result = db.DB.Save(&user)
		if result.Error != nil {
			fmt.Printf("Error promoting user to admin: %v\n", result.Error)
			return
		}

		fmt.Printf("Successfully promoted %s (%s) to admin\n", user.Username, user.Email)
	},
}

var removeAdminCmd = &cobra.Command{
	Use:   "remove-admin [email_or_username]",
	Short: "Remove admin privileges from a user",
	Long:  `Demote an admin user back to regular user role by their email address or username.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		identifier := args[0]
		
		var user types.User
		result := db.DB.Where("email = ? OR username = ?", identifier, identifier).First(&user)
		if result.Error != nil {
			fmt.Printf("Error: User with email/username '%s' not found\n", identifier)
			return
		}

		if user.Role == types.UserRoleUser {
			fmt.Printf("User %s (%s) is already a regular user\n", user.Username, user.Email)
			return
		}

		if user.Role == types.UserRoleSuper {
			fmt.Printf("Cannot demote super admin %s (%s). Use a different method to modify super admin privileges.\n", user.Username, user.Email)
			return
		}

		user.Role = types.UserRoleUser
		result = db.DB.Save(&user)
		if result.Error != nil {
			fmt.Printf("Error demoting user: %v\n", result.Error)
			return
		}

		fmt.Printf("Successfully removed admin privileges from %s (%s)\n", user.Username, user.Email)
	},
}

var createUserCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new user interactively",
	Long:  `Create a new user with interactive prompts for user information.`,
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		// Get email
		fmt.Print("Enter email: ")
		email, _ := reader.ReadString('\n')
		email = strings.TrimSpace(email)
		if email == "" {
			fmt.Println("Error: Email is required")
			return
		}

		// Check if email already exists
		var existingUser types.User
		if result := db.DB.Where("email = ?", email).First(&existingUser); result.Error == nil {
			fmt.Printf("Error: User with email '%s' already exists\n", email)
			return
		}

		// Get username
		fmt.Print("Enter username: ")
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)
		if username == "" {
			fmt.Println("Error: Username is required")
			return
		}

		// Check if username already exists
		if result := db.DB.Where("username = ?", username).First(&existingUser); result.Error == nil {
			fmt.Printf("Error: User with username '%s' already exists\n", username)
			return
		}

		// Get password
		fmt.Print("Enter password: ")
		passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Printf("Error reading password: %v\n", err)
			return
		}
		fmt.Println() // New line after password input
		password := strings.TrimSpace(string(passwordBytes))
		if password == "" {
			fmt.Println("Error: Password is required")
			return
		}

		// Get first name
		fmt.Print("Enter first name (optional): ")
		firstName, _ := reader.ReadString('\n')
		firstName = strings.TrimSpace(firstName)

		// Get last name
		fmt.Print("Enter last name (optional): ")
		lastName, _ := reader.ReadString('\n')
		lastName = strings.TrimSpace(lastName)

		// Ask if user should be admin
		fmt.Print("Make this user an admin? (y/N): ")
		adminInput, _ := reader.ReadString('\n')
		adminInput = strings.TrimSpace(strings.ToLower(adminInput))
		
		role := types.UserRoleUser
		if adminInput == "y" || adminInput == "yes" {
			role = types.UserRoleAdmin
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Printf("Error hashing password: %v\n", err)
			return
		}

		// Create the user
		user := types.User{
			UUID:         uuid.New(),
			Email:        email,
			Username:     username,
			PasswordHash: string(hashedPassword),
			FirstName:    firstName,
			LastName:     lastName,
			Role:         role,
		}

		result := db.DB.Create(&user)
		if result.Error != nil {
			fmt.Printf("Error creating user: %v\n", result.Error)
			return
		}

		roleStr := "user"
		if role == types.UserRoleAdmin {
			roleStr = "admin"
		}

		fmt.Printf("Successfully created %s user: %s (%s)\n", roleStr, username, email)
	},
}

var listUsersCmd = &cobra.Command{
	Use:   "list",
	Short: "List all users",
	Long:  `List all users in the system with their roles and basic information.`,
	Run: func(cmd *cobra.Command, args []string) {
		var users []types.User
		result := db.DB.Find(&users)
		if result.Error != nil {
			fmt.Printf("Error fetching users: %v\n", result.Error)
			return
		}

		if len(users) == 0 {
			fmt.Println("No users found")
			return
		}

		fmt.Printf("%-5s %-20s %-30s %-15s %-12s %s\n", "ID", "Username", "Email", "Role", "2FA", "Created")
		fmt.Println(strings.Repeat("-", 100))

		for _, user := range users {
			twoFA := "disabled"
			if user.TwoFactorEnabled {
				twoFA = "enabled"
			}
			fmt.Printf("%-5d %-20s %-30s %-15s %-12s %s\n", 
				user.ID, 
				user.Username, 
				user.Email, 
				user.Role,
				twoFA,
				user.CreatedAt.Format("2006-01-02"),
			)
		}
	},
}

func init() {
	userCmd.AddCommand(makeAdminCmd)
	userCmd.AddCommand(removeAdminCmd)
	userCmd.AddCommand(createUserCmd)
	userCmd.AddCommand(listUsersCmd)
}