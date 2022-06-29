import { ChatContact, ChatMessage, ChatUnchecked } from '@/socket/ws.types';
import md5 from 'md5';
import { Cookies } from '@/data/cookies';

export enum StoragesKey {
  Partners = 'Partners',
  ChatInputMessage = 'ChatInputMessage',
  ChatMessageList = 'ChatMessageList',
  ChatCurrentContact = 'ChatCurrentContact',
  ChatUncheckData = 'ChatUncheckData',
  LoginStatus = 'LoginStatus',
}

export class Storages {
  static setLoginStatus(value: string | undefined) {
    if (value === undefined) {
      localStorage.removeItem(md5(StoragesKey.LoginStatus));
    } else {
      localStorage.setItem(md5(StoragesKey.LoginStatus), md5(value));
    }
  }

  /**
   *
   * @param value Partner list
   */
  static setPartners(value: ChatContact[]) {
    localStorage.setItem(
      md5(StoragesKey.Partners + Cookies.getUserInfo()?.uid),
      JSON.stringify(value),
    );
  }

  static getPartners(): ChatContact[] {
    const value = localStorage.getItem(
      md5(StoragesKey.Partners + Cookies.getUserInfo()?.uid),
    );
    if (value !== null) {
      return JSON.parse(value);
    }
    return [];
  }

  /**
   *
   * @param key Partner uid
   * @param value message content
   */
  static setChatInputMessage(key: string, value: string) {
    localStorage.setItem(
      md5(StoragesKey.ChatInputMessage + Cookies.getUserInfo()?.uid + key),
      JSON.stringify(value),
    );
  }

  /**
   *
   * @param key Partner uid
   */
  static getChatInputMessage(key): string | null {
    const value = localStorage.getItem(
      md5(StoragesKey.ChatInputMessage + Cookies.getUserInfo()?.uid + key),
    );

    return value;
  }

  /**
   *
   * @param key Partner uid
   * @param value message content list
   */
  static setChatMessage(key: string, value: ChatMessage[]) {
    localStorage.setItem(
      md5(StoragesKey.ChatMessageList + Cookies.getUserInfo()?.uid + key),
      JSON.stringify(value),
    );
  }

  /**
   *
   * @param key Partner uid
   */
  static getChatMessage(key): [ChatMessage] {
    const value = localStorage.getItem(
      md5(StoragesKey.ChatMessageList + Cookies.getUserInfo()?.uid + key),
    );
    if (value !== null) {
      return JSON.parse(value);
    }
    // @ts-ignore
    return [];
  }

  static setCurrentContact(key: string, value: ChatContact) {
    localStorage.setItem(
      md5(StoragesKey.ChatCurrentContact + key),
      JSON.stringify(value),
    );
  }

  static getCurrentContact(key: string): ChatContact | undefined {
    const value = localStorage.getItem(
      md5(StoragesKey.ChatCurrentContact + key),
    );
    if (value) {
      return JSON.parse(value);
    }
    return undefined;
  }

  static setChatUncheckData(value: ChatUnchecked[]) {
    localStorage.setItem(
      md5(StoragesKey.ChatUncheckData + Cookies.getUserInfo()?.uid),
      JSON.stringify(value),
    );
  }

  static getChatUncheckData(): [ChatUnchecked] {
    const value = localStorage.getItem(
      md5(StoragesKey.ChatMessageList + Cookies.getUserInfo()?.uid),
    );
    if (value !== null) {
      return JSON.parse(value);
    }
    // @ts-ignore
    return [];
  }
}
