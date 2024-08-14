namespace CaseConverter {
  const isArray = (arr: any): arr is any[] => {
    return Array.isArray(arr);
  };

  const isObject = (obj: any): obj is Record<string, any> => {
    return obj === Object(obj) && !isArray(obj) && typeof obj !== "function";
  };

  const _toCamel = (value: string): string => {
    return value.replace(/([-_][a-z])/gi, ($1) => {
      return $1.toUpperCase().replace("-", "").replace("_", "");
    });
  };

  const _toSnake = (value: string): string => {
    return value.replace(/[A-Z]/g, (letter) => `_${letter.toLowerCase()}`);
  };

  export const toCamel = (obj: any): any => {
    if (isObject(obj)) {
      const n: Record<string, any> = {};

      Object.keys(obj).forEach((k) => {
        n[_toCamel(k)] = toCamel(obj[k]);
      });

      return n;
    } else if (isArray(obj)) {
      return obj.map((i) => {
        return toCamel(i);
      });
    }

    return obj;
  };

  export const toSnake = (obj: any): any => {
    if (isObject(obj)) {
      const n: Record<string, any> = {};

      Object.keys(obj).forEach((k) => {
        n[_toSnake(k)] = toSnake(obj[k]);
      });

      return n;
    } else if (isArray(obj)) {
      return obj.map((i) => {
        return toSnake(i);
      });
    }

    return obj;
  };
}

export default CaseConverter;
