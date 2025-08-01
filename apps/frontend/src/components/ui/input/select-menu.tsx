import { forwardRef, SelectHTMLAttributes, ReactNode } from "react";
import { InputError } from "./input-error";

interface SelectOption {
  value: string;
  label: string;
  disabled?: boolean;
}

interface SelectMenuProps
  extends Omit<SelectHTMLAttributes<HTMLSelectElement>, "className"> {
  options: SelectOption[];
  error?: string;
  className?: string;
  selectClassName?: string;
  placeholder?: string;
}

export const SelectMenu = forwardRef<HTMLSelectElement, SelectMenuProps>(
  (
    {
      options,
      error,
      className = "",
      selectClassName = "",
      placeholder,
      ...props
    },
    ref
  ) => {
    return (
      <div className={className}>
        <select
          ref={ref}
          className={`
            w-full bg-neutral-900/50 border rounded-md px-3 py-2 text-white 
            focus:outline-none focus:ring-2 transition-colors duration-200
            ${
              error
                ? "border-red-500 focus:ring-red-500 focus:border-red-500"
                : "border-neutral-600 focus:ring-violet-500 focus:border-violet-500"
            }
            ${selectClassName}
          `}
          {...props}
        >
          {placeholder && (
            <option value="" disabled>
              {placeholder}
            </option>
          )}
          {options.map((option) => (
            <option
              key={option.value}
              value={option.value}
              disabled={option.disabled}
              className="bg-neutral-800 text-white"
            >
              {option.label}
            </option>
          ))}
        </select>
        <InputError>{error}</InputError>
      </div>
    );
  }
);

SelectMenu.displayName = "SelectMenu";
