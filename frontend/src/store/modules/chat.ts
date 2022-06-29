import {
  action,
  computed,
  makeObservable,
  observable,
  runInAction,
} from 'mobx';
import ws from '@/socket/ws';
import {
  ChannelTypeEnum,
  ChatLogin,
  ChatContact,
  ChatMessage,
  ChatUnchecked,
} from '@/socket/ws.types';
import { Storages } from '@/data/storages';
import { Cookies } from '@/data/cookies';
import { getChatHistory } from '@/api/modules';

class ChatStore {
  @observable show: boolean = false;
  @observable isLogin: boolean = false;
  @observable chatAuth: string = '';
  @observable private _partnersList: ChatContact[] = [];
  @observable private _chatUncheckedData: ChatUnchecked[] = [];
  @observable private _chatMessages: ChatMessage[] = [];
  @observable private _currentContact?: ChatContact;
  @observable isAutoScroll: boolean = true;
  @observable contactSearchText: string = '';

  @computed get isChatAuth() {
    return this.chatAuth === 'ok';
  }

  @computed get partnersList() {
    const _list =
      this._partnersList.length === 0
        ? Storages.getPartners()
        : this._partnersList;

    return this.contactSearchText.length === 0
      ? _list
      : _list.filter((value) => {
          let str = `\S*${this.contactSearchText}\S*`;
          let reg = new RegExp(str, 'i');
          return reg.test(value.name);
        });
  }

  @computed get chatUncheckedData() {
    return this._chatUncheckedData.length === 0
      ? Storages.getChatUncheckData()
      : this._chatUncheckedData;
  }

  @computed get chatMessages() {
    return this._chatMessages.length === 0
      ? Storages.getChatMessage(this.currentContact?.uid)
      : this._chatMessages;
  }

  @computed get currentContact(): ChatContact | undefined {
    const userInfo = Cookies.getUserInfo();
    if (!userInfo) return undefined;
    const contact =
      this._currentContact === undefined
        ? Storages.getCurrentContact(userInfo.uid)
          ? Storages.getCurrentContact(userInfo.uid)
          : undefined
        : this._currentContact;
    return contact;
  }

  constructor() {
    makeObservable(this);
    ws.addEventListener(ChannelTypeEnum.Login, (val: ChatLogin) => {
      this.SET_CHAT_AUTH(val.msg);
    });
    ws.addEventListener(ChannelTypeEnum.ChatPartners, (val: ChatContact[]) => {
      const userInfo = Cookies.getUserInfo();
      if (userInfo) {
        const currentContact = Storages.getCurrentContact(userInfo.uid);
        if (!currentContact) {
          this.SET_CURRENT_CONTACT(val[0]);
          Storages.setCurrentContact(userInfo.uid, val[0]);
        }
      }
      this.SET_PARTNER_DATA(val);
    });

    ws.addEventListener(ChannelTypeEnum.ChatHistory, (val: ChatMessage[]) => {
      this.SET_MESSAGES(val.reverse());
    });

    ws.addEventListener(ChannelTypeEnum.ChatReceive, (val: ChatMessage) => {
      const current = this.currentContact;
      if (current?.uid === val.sender || current?.uid === val.recipient) {
        /// it is current data
        this.SET_MESSAGE(val);
      } else {
        const value = Storages.getChatMessage(val.recipient);
        if (val.recipient) {
          Storages.setChatMessage(val.recipient, [...value, val]);
        }
      }
    });
    ws.addEventListener(
      ChannelTypeEnum.ChatUncheckedList,
      (val: ChatUnchecked[]) => {
        this.SET_UN_CHECK_DATA(val);
      },
    );
  }

  @action
  SET_CONTACT_SEARCH(value: string) {
    runInAction(() => {
      this.contactSearchText = value;
    });
  }

  @action
  SET_AUTO_SCROLL(val: boolean) {
    runInAction(() => {
      this.isAutoScroll = val;
    });
  }

  @action
  SET_CURRENT_CONTACT(val: ChatContact) {
    runInAction(() => {
      this._currentContact = val;
    });
    const userInfo = Cookies.getUserInfo();
    if (userInfo) {
      Storages.setCurrentContact(userInfo.uid, val);
      const messages = Storages.getChatMessage(val.uid);
      this.SET_MESSAGES(messages);
    }
  }

  @action
  SET_SHOW(_show: boolean) {
    runInAction(() => {
      this.show = _show;
    });
  }

  @action
  SET_Login(_login: boolean) {
    runInAction(() => {
      this.isLogin = _login;
    });
  }

  @action
  SET_CHAT_AUTH(_auth: string) {
    runInAction(() => {
      this.chatAuth = _auth;
    });
  }

  removeDuplicate = (data: ChatContact[]) => {
    let len = data.length;
    for (let i = 0; i < len; i++) {
      for (let j = i + 1; j < len; j++) {
        if (data[i].uid === data[j].uid) {
          data.splice(j, 1);
          len--; // 减少循环次数提高性能
          j--; // 保证j的值自加后不变
        }
      }
    }
    return data;
  };

  /**
   * if you update local and store data must  call this function
   * @param value
   * @constructor
   */
  @action
  SET_PARTNER_DATA(value: ChatContact[]) {
    let data = Storages.getPartners();
    data.push(...value);
    data = this.removeDuplicate(data);
    Storages.setPartners(data);
    runInAction(() => {
      this._partnersList = data;
    });
  }

  @action
  SET_MESSAGE(val: ChatMessage, save: boolean = true) {
    this.SET_MESSAGES([...this.chatMessages, val], save);
  }

  @action
  SET_INSERT_MESSAGES(val: ChatMessage[], save: boolean = true) {
    this.SET_MESSAGES([...val, ...this.chatMessages], save);
  }

  @action
  SET_MESSAGES(val: ChatMessage[], save: boolean = true) {
    if (this.currentContact && save) {
      Storages.setChatMessage(this.currentContact.uid, val);
    }
    runInAction(() => {
      this._chatMessages = val;
    });
  }

  @action
  SET_UN_CHECK_DATA(val: ChatUnchecked[]) {
    Storages.setChatUncheckData(val);
    runInAction(() => {
      this._chatUncheckedData = val;
    });
  }

  @action
  async getChatHistory() {
    const _message = this.chatMessages;
    getChatHistory({
      partner: this.currentContact?.uid,
      before: _message.length == 0 ? null : this.chatMessages[0].timestamp,
      limit: 20,
    })
      .then((value) => {
        this.SET_AUTO_SCROLL(false);
        runInAction(() => {
          this.SET_INSERT_MESSAGES(
            // @ts-ignore
            (value.data as ChatMessage[]).reverse(),
            false,
          );
        });
      })
      .catch((e) => {
        console.log(':=======E===', e);
      });
  }
}

const chatStore = new ChatStore();

export default chatStore;
