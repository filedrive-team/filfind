import Header from '@/components/Header';
import Footer from '@/components/Footer';
import Filter from '@/components/ClientList/Filter';
import List from '@/components/ClientList/List';
import styles from './style.module.scss';
import { useEffect } from 'react';
import headerStore from '@/store/modules/header';
import clientListStore from '@/store/modules/clientList';
const ClientList = () => {
  headerStore.SET_ACTIVEINDEX(1);
  useEffect(() => {
    clientListStore._getClients(clientListStore.params);
  }, []);
  return (
    <div className={styles.clientList}>
      <div className={styles.bgWrap}></div>
      <Header></Header>
      <Filter></Filter>
      <List></List>
      <Footer></Footer>
    </div>
  );
};

export default ClientList;
