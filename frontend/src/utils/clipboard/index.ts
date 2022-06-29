import ClipboardJS from 'clipboard';

export function copyText(text: string, container?: Element): Promise<any> {
  return new Promise(function (resolve, reject) {
    const fakeElement = document.createElement('button');
    const clipboard = new ClipboardJS(fakeElement, {
      text: function () {
        return text;
      },
      action: function () {
        return 'copy';
      },
      container: typeof container === 'object' ? container : document.body,
    });

    clipboard.on('success', function (e) {
      clipboard.destroy();
      resolve(e);
    });

    clipboard.on('error', function (e) {
      clipboard.destroy();
      reject(e);
    });

    fakeElement.click();
  });
}
