import styles from './List.module.scss';
import { Table } from 'antd';
import clientStore from '@/store/modules/client';
import { observer } from 'mobx-react';
import { bytesToSize, filToSize, formatRate } from '@/utils';
import dayjs from 'dayjs';

const List = () => {
  const columns = [
    {
      title: 'No',
      key: 'No',
      render: (text, record, index) => (
        <div
          style={{
            width: '25px',
            height: ' 25px',
            opacity: '0.6',
            border: '1px solid #475569',
            borderRadius: '50%',
            lineHeight: '25px',
            textAlign: 'center',
          }}
        >
          {index + 1}
        </div>
      ),
    },
    {
      title: 'Miner ID',
      dataIndex: 'provider',
      key: 'provider',
      render: (text, record) => {
        const { owner } = record;
        return (
          <div onClick={() => skipToSp(owner)} className="listLink">
            {text}
          </div>
        );
      },
    },
    {
      title: 'Country',
      dataIndex: 'iso_code',
      key: 'iso_code',
      sorter: {
        compare: (a, b) => {
          const _code_a = a.iso_code ?? '';
          const _code_b = b.iso_code ?? '';
          return _code_a.toLocaleUpperCase().localeCompare(_code_b);
        },
        multiple: 1,
      },
      render: (text) => {
        return (
          <div
            style={{
              width: '100px',
            }}
          >
            {text === '' ? '-' : text}
          </div>
        );
      },
    },
    {
      title: 'Region',
      dataIndex: 'region',
      key: 'region',
      render: (text) => {
        return (
          <div
            style={{
              width: '100px',
            }}
          >
            {text === '' ? '-' : text}
          </div>
        );
      },
    },
    {
      title: 'Total Storage Deals',
      dataIndex: 'storage_deals',
      key: 'storage_deals',
      render: (text) => {
        return (
          <div
            style={{
              width: '120px',
            }}
          >
            {text === '' ? '-' : text}
          </div>
        );
      },
    },
    {
      title: 'Storage Capacity',
      dataIndex: 'storage_capacity',
      key: 'storage_capacity',
      sorter: {
        compare: (a, b) =>
          Number(a.storage_capacity ?? 0) - Number(b.storage_capacity ?? 0),
        multiple: 2,
      },
      render: (text, record) => {
        const _format = bytesToSize(text);
        return (
          <div
            style={{
              width: '120px',
            }}
          >
            {_format}
          </div>
        );
      },
    },
    {
      title: 'Total DataCap Allocation',
      dataIndex: 'data_cap',
      key: 'data_cap',
      sorter: {
        compare: (a, b) => Number(a.data_cap ?? 0) - Number(b.data_cap ?? 0),
        multiple: 3,
      },
      render: (text, record) => {
        const _format = bytesToSize(text);
        return (
          <div
            style={{
              width: '120px',
            }}
          >
            {_format}
          </div>
        );
      },
    },
    {
      title: 'Avg.Price for Non-verified Deals',
      dataIndex: 'avg_price',
      key: 'avg_price',
      sorter: {
        compare: (a, b) => Number(a.avg_price) - Number(b.avg_price),
        multiple: 4,
      },
      render: (text, record) => {
        const _format = filToSize(text);
        return (
          <div
            style={{
              width: '120px',
            }}
          >
            {_format}
          </div>
        );
      },
    },
    {
      title: 'Avg.Price for Verified Deals',
      dataIndex: 'avg_verified_price',
      key: 'avg_verified_price',
      sorter: {
        compare: (a, b) =>
          Number(a.avg_verified_price) - Number(b.avg_verified_price),
        multiple: 5,
      },
      render: (text, record) => {
        const _format = filToSize(text);
        return (
          <div
            style={{
              width: '120px',
            }}
          >
            {_format}
          </div>
        );
      },
    },
    {
      title: 'Time of 1st-Deal',
      dataIndex: 'first_deal_time',
      key: 'first_deal_time',
      render: (text, record) => {
        const _format =
          text === 0 ? '-' : dayjs.unix(text).format('YYYY-MM-DD');
        return (
          <div
            style={{
              wordWrap: 'break-word',
              wordBreak: 'break-word',
              width: '100px',
            }}
          >
            {_format}
          </div>
        );
      },
    },
    {
      title: 'Deal Success Rate',
      dataIndex: 'storage_success_rate',
      key: 'storage_success_rate',
      sorter: {
        compare: (a, b) =>
          Number(a.storage_success_rate) - Number(b.storage_success_rate),
        multiple: 6,
      },
      render: (text, record) => {
        return (
          <div
            style={{
              width: '100px',
            }}
          >
            {formatRate(text)}
          </div>
        );
      },
    },
    {
      title: 'Retrieval Success Rate',
      dataIndex: 'retrieval_success_rate',
      key: 'retrieval_success_rate',
      sorter: {
        compare: (a, b) =>
          Number(a.retrieval_success_rate) - Number(b.retrieval_success_rate),
        multiple: 7,
      },
      render: (text, record) => {
        return <div>{formatRate(text)}</div>;
      },
    },
  ];
  const skipToSp = (id) => {
    window.open(`/storageProvide/${id}`);
  };
  return (
    <div className={styles.list}>
      <div className={styles.title}>History Deal Record</div>
      <Table
        columns={columns}
        dataSource={clientStore.historyDealStatsList}
        rowKey={(record) => record.provider}
        pagination={{
          total: clientStore.historyDealStatsTotal,
          onChange: (page, pageSize) => {
            clientStore.SET_HISTORYDEALSTATSPARAM({
              ...clientStore.historyDealStatsParam,
              page: page,
              page_size: pageSize,
            });
            clientStore._getClientHistoryDealStats(
              clientStore.historyDealStatsParam,
            );
          },
        }}
      />
    </div>
  );
};

export default observer(List);
