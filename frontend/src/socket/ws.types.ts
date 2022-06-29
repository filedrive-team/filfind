export declare class MessageData {
  channel: ChannelTypeEnum;
  subscribe: boolean;
  body?: any;
}

export declare class SendData {
  channel: ChannelTypeEnum;
  subscribe: boolean;
  body?: any;
}

export enum MessageTypeEnum {
  Subscriptions = 'SUBSCRIBE',
  Snapshot = 'snapshot',
  Update = 'update',
  Pong = 'pong',
}

export enum ChannelTypeEnum {
  Echo = 'echo',
  Login = 'login',
  Logout = 'logout',
  Auth = 'auth',
  Ping = 'ping',
  ChatSend = 'chat_send',
  ChatReceive = 'chat_receive',
  ChatPartners = 'chat_partners',
  ChatPartnersStatus = 'chat_partners_status',
  ChatHistory = 'chat_history',
  ChatCheck = 'chat_check',
  ChatCheckedStatus = 'chat_checked_status',
  ChatUncheckedTotal = 'chat_unchecked_total',
  ChatUncheckedList = 'chat_unchecked_list',
}

export declare class ChatLogin {
  msg: string;
}

export declare class ChatSend {
  recipient: string;
  type: string;
  content: string;
}

export declare class ChatMessage {
  recipient?: string;
  sender?: string;
  type: number;
  content: string;
  timestamp: number;
  checked: boolean;
}

export declare class ChatContact {
  uid: string;
  type: string;
  name: string;
  address_id: string;
  avatar: string;
  location: string;
  contact_email: string;
  slack?: string;
  github?: string;
  twitter?: string;
  description?: string;
  unCheckCount?: number;
}

export declare class ChatPartnersStatus {
  uid: string;
  online: boolean;
}

export declare class ChatHistorySend {
  partner: string;
  before: number;
  limit: number;
}

export declare class ChatUncheckedTotal {
  total: number;
}

export declare class ChatCheck {
  partner: string;
}

export declare class ChatUnchecked {
  partner: string;
  number: number;
}
