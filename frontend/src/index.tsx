import React from 'react';
import ReactDOM from 'react-dom';
import './index.scss';
import './assets/style/var.scss';
import { Provider } from 'mobx-react';
import * as stores from './store';
import 'antd/dist/antd.css';
import App from './App';
import '@/socket/ws';

ReactDOM.render(
  <React.StrictMode>
    <Provider {...stores}>
      <App />
    </Provider>
  </React.StrictMode>,
  document.getElementById('root'),
);
