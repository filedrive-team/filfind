import EventEmitter from '@/utils/event';
import { Subject } from 'rxjs';
import { ChannelTypeEnum, MessageData, SendData } from '@/socket/ws.types';
import { Cookies } from '@/data/cookies';

class CWebSocket extends EventEmitter {
  readonly timeout = 2000;
  readonly retryInterval = 2000;
  lockReconnect = false;
  wsUrl: string;
  timer: NodeJS.Timeout | undefined;
  serverTimer: NodeJS.Timeout | undefined;
  ws: WebSocket | undefined;
  openCount = 0;
  token = '';
  subject = new Subject<MessageData>();
  subscribeSet = new Set<string>();
  reconnectTimes = 0;

  constructor(url: string) {
    super();
    this.wsUrl = url;
    this.createWebSocket(url);
  }

  private createWebSocket(url: string) {
    this.ws = new WebSocket(url);
    this.ws.onopen = this.onopen;
    this.ws.onmessage = this.onmessage;
    this.ws.onclose = (val) => {
      this.reconnectWebSocket();
    };
    this.ws.onerror = (val) => {
      if (this.reconnectTimes >= 0) {
        this.reconnectTimes = -2;
      } else {
        this.reconnectTimes += 1;
      }
      this.reconnectWebSocket();
    };
  }

  reconnectWebSocket() {
    if (this.lockReconnect) return;
    this.lockReconnect = true;
    setTimeout(() => {
      this.createWebSocket(this.wsUrl);
      this.lockReconnect = false;
    }, this.retryInterval);
  }

  send(data: SendData) {
    if (this.ws?.readyState !== 1) {
      setTimeout(() => this.send(data), 100);
      return;
    }
    this.ws.send(JSON.stringify(data));
  }

  setToken(token: string) {
    this.token = token;
  }

  auth() {
    if (Cookies.getUserInfo() !== null) return;
    this.send({
      channel: ChannelTypeEnum.Login,
      subscribe: false,
      body: {
        token: Cookies.getUserInfo()?.access_token,
      },
    });
  }

  ping() {
    this.send({ channel: ChannelTypeEnum.Ping, subscribe: true });
  }

  // subscribe(channels: any[]) {
  //   channels.forEach((e) => {
  //     this.subscribeSet.add(JSON.stringify(e));
  //   });
  //   this.send({ type: 'subscribe', channels });
  // }
  //
  // unsubscribe(channels: any[]) {
  //   channels.forEach((e) => {
  //     this.subscribeSet.delete(JSON.stringify(e));
  //   });
  //   this.send({ subscribe: true, channels });
  // }

  onopen = (val) => {
    if (this.openCount !== 0) {
      this.auth();
    }
    this.openCount += 1;
    this.startHeartCheck();
    this.dispatchEvent('open');
  };

  onmessage = (e: MessageEvent) => {
    try {
      this.resetHeartCheck();
      const _data = JSON.parse(e.data);
      if (typeof _data === 'number') return;
      const data = _data as MessageData;
      const { channel, body } = data;

      if (channel !== undefined) {
        this.dispatchEvent(channel, body);
      }
    } catch (error) {
      console.log('======onmessage===error====', e.data, '\n', error, '=====');
    }
  };

  private startHeartCheck() {
    this.timer = setTimeout(() => {
      // this.ping();
    }, this.timeout);
  }

  private clearHeartCheck() {
    if (this.timer) {
      clearTimeout(this.timer);
    }
  }

  private resetHeartCheck() {
    this.clearHeartCheck();
    this.startHeartCheck();
  }
}

const ws = new CWebSocket(process.env['REACT_APP_BASE_WS']!);

export default ws;
