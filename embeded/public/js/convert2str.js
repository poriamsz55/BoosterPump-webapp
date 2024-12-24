
export function converterToString(value) {
    if (value===undefined) {
        return "";
    }

    let converterStr = '';
    switch (value.toString()) {
        case "0":
            converterStr = 'بدون تبدیل'
            break;
        case "1":
            converterStr = 'تبدیل در دهش'
            break;
        case "2":
            converterStr = 'تبدیل دو طرفه'
            break;

        default:
            break;
    }

    return converterStr;
}

export function filterToString(value) {
    let filterStr = '';
    if (value) {
        filterStr = 'دارد';
    }else{
        filterStr = 'ندارد';
    }
    return filterStr;
}