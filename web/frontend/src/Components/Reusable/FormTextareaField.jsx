import React from "react";
function FormTextareaField({
  label,
  name,
  placeholder,
  value,
  errorText,
  helpText,
  onChange,
  isRequired,
  maxWidth,
  rows = 2,
  disabled = false,
}) {
  let classNameText = "textarea";
  if (errorText) {
    classNameText = "textarea is-danger";
  }
  return (
    <div class="field pb-4">
      <label class="label">{label}</label>
      <div class="control">
        <textarea
          className={classNameText}
          name={name}
          placeholder={placeholder}
          value={value}
          onChange={onChange}
          style={{ maxWidth: maxWidth }}
          rows={rows}
          disabled={disabled}
          autoComplete="off"
        />
      </div>
      {errorText && <p class="help is-danger">{errorText}</p>}
      {helpText && <p class="help">{helpText}</p>}
    </div>
  );
}

export default FormTextareaField;
