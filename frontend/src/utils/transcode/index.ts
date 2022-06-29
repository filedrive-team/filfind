export function strToHex(str: string): string {
  const arr: any = [];
  let length = str.length;
  for (let i = 0; i < length; i++) {
    const hex = str.charCodeAt(i).toString(16);
    arr[i] = hex.length === 1 ? '0' + hex : hex;
  }
  return arr.join('');
}
