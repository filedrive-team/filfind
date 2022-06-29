const PullRefresh = (top: string, content: string) => {
  let clientX = 0;
  let clientY = 0;
  let PMoveX = 0;
  let PstartX = 0;
  let PstartY = 0;
  let PMoveY = 0;
  let PendX = 0;
  let PendY = 0;

  var flag = false;
  document.onmousedown = function (ev) {
    flag = true;
    PstartX = ev.pageX;
    PstartY = ev.pageY;
    document.onmousemove = function (ev) {
      PMoveX = ev.pageX;
      PMoveY = ev.pageY;
      if (flag) {
        var resutl = getpostion(PMoveY, PstartY);
        switch (resutl) {
          case 0:
            break;
          case 1:
            break;
          case 2:
            if (PMoveY - PstartY > 0) {
              if (PMoveY - PstartY >= 100) {
                // @ts-ignore
                document.getElementById(id).style.marginTop = '100px';
              } else {
                // @ts-ignore
                document.getElementById(id).style.marginTop =
                  PMoveY - PstartY + 'px';
              }
              // @ts-ignore
              document.getElementById(id).style.display = 'block';
            }
            break;
        }
      }
    };
    document.onmouseup = function (ev) {
      console.log('====onmouseup====');

      flag = false;
      PendX = ev.pageX;
      PendY = ev.pageY;
      var resutl = getpostion(PMoveY, PstartY);
      switch (resutl) {
        case 0:
          break;
        case 1:
          break;
        case 2:
          break;
      }
    };

    function getpostion(PMoveY, PstartY) {
      if (PMoveY - PstartY == 0) {
        return 0;
      }
      if (PMoveY - PstartY < 0) {
        return 1;
      }
      if (PMoveY - PstartY > 0) {
        return 2;
      }
    }
  };
};

export { PullRefresh };
