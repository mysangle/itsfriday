
export function toSnakeCase(obj: Record<string, any>): Record<string, any> {
  const result: Record<string, any> = {};
  for (const key in obj) {
    const snakeKey = key
      .replace(/([A-Z])/g, "_$1") // 대문자 앞에 '_' 추가
      .toLowerCase();
    result[snakeKey] = obj[key];
  }
  return result;
}

export function toCamelCase<T extends Record<string, any>>(obj: T): any {
  const result: Record<string, any> = {};
  for (const key in obj) {
    const camelKey = key.replace(/_([a-z])/g, (_, char) => char.toUpperCase());
    result[camelKey] = obj[key];
  }
  return result;
}
