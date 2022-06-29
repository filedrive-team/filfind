class EventEmitter {
  listeners = {};

  addEventListener(type: string, callback: Function) {
    if (!(type in this.listeners)) {
      this.listeners[type] = [];
    }
    this.listeners[type].push(callback);
  }

  removeEventListener(type: string, callback: Function) {
    if (!(type in this.listeners)) {
      return;
    }
    const stack = this.listeners[type];
    for (let i = 0, l = stack.length; i < l; i++) {
      if (stack[i] === callback) {
        stack.splice(i, 1);
        this.removeEventListener(type, callback);
        return;
      }
    }
  }

  dispatchEvent(type: string, ...args: any[]) {
    if (!(type in this.listeners)) {
      return;
    }
    const stack = this.listeners[type];
    for (let i = 0, l = stack.length; i < l; i++) {
      stack[i].call(this, ...args);
    }
  }
}

export default EventEmitter;
