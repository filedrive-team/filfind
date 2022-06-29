import Header from '@/components/Header';
import Footer from '@/components/Footer';
import Filter from '@/components/Home/Filter';
import SearchCondition from '@/components/Home/SearchCondition';
import List from '@/components/Home/List';
import styles from './style.module.scss';
import bgLeft from '@/assets/images/home/bg_left.png';
import bgRight from '@/assets/images/home/bg_right.png';
import { useEffect } from 'react';
import homeStore from '@/store/modules/home';
import headerStore from '@/store/modules/header';
import { Cookies } from '@/data/cookies';
const Home = () => {
  headerStore.SET_ACTIVEINDEX(0);
  useEffect(() => {
    const _user = Cookies.getUserInfo();
    console.log(_user, '_user');

    homeStore._getProviders(homeStore.params);
  }, []);
  return (
    <div className={styles.home}>
      <Header></Header>
      <div className={styles.bgLeft}>
        <img src={bgLeft} alt="" />
      </div>
      <div className={styles.bgRight}>
        <img src={bgRight} alt="" />
      </div>
      <div className={styles.top}>
        <div className={styles.filterWrap}>
          <Filter></Filter>
        </div>
        <div className={styles.searchConditionWrap}>
          <SearchCondition></SearchCondition>
        </div>
      </div>
      <div className={styles.listWrap}>
        <List></List>
      </div>
      <Footer></Footer>
    </div>
  );
};

export default Home;
