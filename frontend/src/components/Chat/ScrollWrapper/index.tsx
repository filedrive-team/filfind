import React, { CSSProperties, MouseEventHandler, ReactNode } from 'react';
import Internal from '@/components/Chat/ScrollWrapper/Internal';
import { ChatContact } from '@/socket/ws.types';
type Callback = (value: boolean) => void;
type IProps = {
  data: Object[];
  scrollToBottom?: boolean;
  children?: ReactNode;
  onEarlier?: MouseEventHandler;
  onIsAutoToScroll?: Callback;
  isAutoToScroll?: boolean;
  me?: ChatContact;
  other?: ChatContact;
  style?: CSSProperties;
};

const ScrollWrapper = (Comp: any) => (props: IProps) => {
  return <Internal component={Comp} {...props} />;
};

export default ScrollWrapper;
