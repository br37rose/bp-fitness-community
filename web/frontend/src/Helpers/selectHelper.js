
export function getSelectedOptions(options, values) {
    try {
        const selectedOptions = options.filter(option => values.includes(option.value));
        return selectedOptions;
    }
    catch(err) {
        return [];
    }
}
