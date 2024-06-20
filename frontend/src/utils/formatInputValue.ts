// function that converts string to float with possible comma sign insted of dot sign
// returns formatted input value and converted float value
export const formatInputValue = (inputValue: string) => {
  let numberValue = 0;
  if (inputValue !== "") {
    if (inputValue.indexOf(",") < -1) {
      numberValue = parseFloat(inputValue);
    } else {
      inputValue = inputValue.replace(",", ".").replace(" ", "");
      numberValue = parseFloat(inputValue);
    }
  }
  inputValue = inputValue.replace(/[^0-9.]/g, "");
  return { inputValue, numberValue };
};
