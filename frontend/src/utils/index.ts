import Decimal from 'decimal.js';

export function desensitization(str, beginLen, endLen) {
  const len = str.length;
  if (beginLen + endLen > len) return str;
  const firstStr = str.substr(0, beginLen);
  const lastStr = str.substr(len - endLen, endLen);
  const tempStr = firstStr + '***' + lastStr;
  return tempStr;
}

export function jwt_decode(value) {
  let string = decodeURIComponent(escape(window.atob(value.split('.')[1])));
  return JSON.parse(string);
}

export function bytesToSize(bytes) {
  if (bytes === '0') return '0 B';
  var k = 1024,
    sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'],
    i = Math.floor(Math.log(bytes) / Math.log(k));
  const _v = bytes / Math.pow(k, i);
  const _d = new Decimal(_v).toDecimalPlaces(2, Decimal.ROUND_DOWN);
  return _d + ' ' + sizes[i];
}

export function formatRate(value) {
  const str = new Decimal(Number(value) * 100).toDecimalPlaces(
    3,
    Decimal.ROUND_DOWN,
  );
  return str + '%';
}

export function filToSize(value: string | null) {
  const _err = ['', null, undefined];
  if (_err.includes(value)) return '-';
  if (value === '0') return '0 FIL';
  const autoFil = new Decimal(value ?? '0').toNumber();
  const k = Math.pow(10, 9);
  const maxFil = 1000 * Math.pow(10, 18);
  if (autoFil > maxFil) return -1;
  const sizes = ['autoFIL', 'nanoFIL', 'FIL'];
  const i = Math.floor(Math.log(autoFil) / Math.log(k));

  const _v = autoFil / Math.pow(k, i);
  const _res = new Decimal(_v).toDecimalPlaces(3, Decimal.ROUND_DOWN);
  return _res + ' ' + sizes[i];
}

export function autoFILToFil(value: string) {
  if (value.length > 18) {
    const start = value.substring(0, value.length - 18);
    let end = value.substring(value.length - 18, value.length);
    let i = end.length - 1;
    while (Number(end[i]) === 0) {
      console.log(end[i], 'end[i]');
      end = end.substring(0, end.length - 1);
      i--;
    }
    const p = end === '' ? '' : '.';
    return start + p + end + 'FIL';
  }
}

export function getWeb3Image(cid: string) {
  if (cid !== '') {
    return `https://dweb.link/ipfs/${cid}`;
  } else {
    return '';
  }
}

export const emailReg =
  /^\w+((-\w+)|(\.\w+))*\@[A-Za-z0-9]+((\.|-)[A-Za-z0-9]+)*\.[A-Za-z0-9]+$/;
