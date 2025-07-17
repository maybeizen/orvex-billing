import { forwardRef, InputHTMLAttributes, ReactNode } from "react";
import { InputError } from "./input-error";

interface CheckboxProps
  extends Omit<InputHTMLAttributes<HTMLInputElement>, "type"> {
  label?: ReactNode;
  error?: string;
  className?: string;
  checkboxClassName?: string;
  labelClassName?: string;
  containerClassName?: string;
}

export const Checkbox = forwardRef<HTMLInputElement, CheckboxProps>(
  (
    {
      label,
      error,
      className = "",
      checkboxClassName = "",
      labelClassName = "",
      containerClassName = "",
      ...props
    },
    ref
  ) => {
    const checked = Boolean(props.checked);

    return (
      <div className={containerClassName}>
        <label
          className={`flex items-center gap-3 cursor-pointer select-none ${className}`}
        >
          <span className="relative flex items-center">
            <input
              ref={ref}
              type="checkbox"
              className="peer appearance-none w-5 h-5 border border-white/20 rounded-md bg-white/5 checked:bg-violet-500 checked:border-violet-500 focus:ring-2 focus:ring-violet-400/40 transition-all duration-200 outline-none"
              {...props}
            />
            <span
              className={`
                pointer-events-none absolute left-0 top-0 w-5 h-5 flex items-center justify-center
                ${error ? "border-red-500 ring-2 ring-red-400/40" : ""}
                ${checkboxClassName}
              `}
            >
              {checked && <i className="fas fa-check text-sm text-white" />}
            </span>
          </span>
          {label && (
            <span
              className={`text-sm text-nowrap text-white ${labelClassName}`}
            >
              {label}
            </span>
          )}
        </label>
        {error && <InputError>{error}</InputError>}
      </div>
    );
  }
);

Checkbox.displayName = "Checkbox";
