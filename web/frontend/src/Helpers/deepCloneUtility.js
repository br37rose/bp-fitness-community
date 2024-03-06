export default function deepClone(obj) {
	if (obj === null || typeof obj !== 'object') {
	// If 'obj' is not an object or is null, return it as is
	return obj;
	}

	if (Array.isArray(obj)) {
	// If 'obj' is an array, create a new array and clone its elements
	const clone = [];
	for (let i = 0; i < obj.length; i++) {
	  clone[i] = deepClone(obj[i]);
	}
	return clone;
	}

	if (obj instanceof Date) {
	// If 'obj' is a Date object, create a new Date object with the same timestamp
	return new Date(obj.getTime());
	}

	if (obj instanceof RegExp) {
	// If 'obj' is a RegExp object, create a new RegExp object with the same pattern and flags
	return new RegExp(obj);
	}

	// If 'obj' is a regular object, create a new object and deep clone its properties
	const clone = {};
	for (let key in obj) {
	if (obj.hasOwnProperty(key)) {
	  clone[key] = deepClone(obj[key]);
	}
	}
	return clone;
}
