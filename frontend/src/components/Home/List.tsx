import styles from './List.module.scss';
import { Table, Tooltip } from 'antd';
import homeStore from '@/store/modules/home';
import { observer } from 'mobx-react';
import { ExclamationCircleOutlined } from '@ant-design/icons';
import { autoFILToFil, bytesToSize, filToSize, formatRate } from '@/utils';
import dayjs from 'dayjs';
const List = () => {
  const columns = [
    {
      title: 'Rank',
      dataIndex: 'rank',
      key: 'rank',
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
      dataIndex: 'address',
      key: 'address',
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
      title: 'SP Name',
      dataIndex: 'name',
      key: 'name',
      render: (text) => {
        return (
          <div
            style={{
              width: '80px',
            }}
          >
            {text === '' ? '-' : text}
          </div>
        );
      },
    },
    {
      title: 'Country',
      dataIndex: 'iso_code',
      key: 'iso_code',
      sorter: true,
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
      title: 'Raw Power',
      dataIndex: 'raw_power',
      key: 'raw_power',
      render: (text, record) => {
        const _format = bytesToSize(text);
        return (
          <div
            style={{
              width: '90px',
            }}
          >
            {_format}
          </div>
        );
      },
    },
    {
      title: 'Adj. Power',
      dataIndex: 'quality_adj_power',
      key: 'quality_adj_power',
      sorter: true,
      render: (text, record) => {
        const _format = bytesToSize(text);
        return (
          <div
            style={{
              width: '90px',
            }}
          >
            {_format}
          </div>
        );
      },
    },
    {
      title: 'Total Storage Deals',
      dataIndex: 'storage_deals',
      key: 'storage_deals',
      sorter: true,
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
      title: 'Price for Non-verified Deals',
      dataIndex: 'price',
      key: 'price',
      sorter: true,
      render: (text, record) => {
        const _format = filToSize(text);
        if (_format === -1) {
          return (
            <div
              style={{
                width: '120px',
              }}
            >
              <Tooltip title={autoFILToFil(text)}>NOT RCMD</Tooltip>
            </div>
          );
        } else {
          return (
            <div
              style={{
                width: '120px',
              }}
            >
              {_format}
            </div>
          );
        }
      },
    },
    {
      title: 'Price for Verified Deals',
      dataIndex: 'verified_price',
      key: 'verified_price',
      sorter: true,
      render: (text, record) => {
        const _format = filToSize(text);
        if (_format === -1) {
          return (
            <div
              style={{
                width: '120px',
              }}
            >
              <Tooltip title={autoFILToFil(text)}>NOT RCMD</Tooltip>
            </div>
          );
        } else {
          return (
            <div
              style={{
                width: '120px',
              }}
            >
              {_format}
            </div>
          );
        }
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
      sorter: true,
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
      sorter: true,
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
      title: 'Special Distinction Features',
      key: 'specialDistinctionFeatures',
      render: (text, record, index) => {
        const { certification, renewable_energy } = record;

        if (certification && renewable_energy) {
          return (
            <div>
              <Tooltip
                title={
                  <div>
                    <div>Renewable Energy:{renewable_energy}</div>
                    <div>Certification:{certification}</div>
                  </div>
                }
              >
                <div className="listTips">2</div>
              </Tooltip>
            </div>
          );
        } else if (certification || renewable_energy) {
          return (
            <div>
              <Tooltip
                title={
                  <div>
                    <div>Renewable Energy:{renewable_energy}</div>
                    <div>Certification:{certification}</div>
                  </div>
                }
              >
                <div className="listTips">1</div>
              </Tooltip>
            </div>
          );
        } else {
          return <div>N/A </div>;
        }
      },
    },
    {
      title: 'Reputation Score',
      dataIndex: 'reputation_score',
      key: 'reputation_score',
      sorter: true,
    },
    {
      title: 'Client Review Score',
      dataIndex: 'review_score',
      key: 'review_score',
      sorter: true,
    },
    {
      title: 'Storage Service Details',
      key: 'storage Service Details',
      render: (text, record, index) => {
        const { available_deals, bandwidth, sealing_speed, parallel_deals } =
          record;
        if (available_deals || bandwidth || sealing_speed || parallel_deals) {
          return (
            <div>
              <Tooltip
                title={
                  <div>
                    <div>
                      {available_deals
                        ? `Arallel Fil + deals:${available_deals}`
                        : ''}
                    </div>
                    <div>
                      {bandwidth ? `Data transfer bandwidth:${bandwidth}` : ''}
                    </div>
                    <div>
                      {sealing_speed ? `Sealing speed:${sealing_speed}` : ''}
                    </div>
                    <div>
                      {parallel_deals
                        ? `Number parallel deals:${parallel_deals}`
                        : ''}
                    </div>
                  </div>
                }
              >
                <ExclamationCircleOutlined className="listTips" />
              </Tooltip>
            </div>
          );
        } else {
          return <div className="disabildTips">N/A</div>;
        }
      },
    },
  ];
  const skipToSp = (id) => {
    window.open(`/storageProvide/${id}`);
  };
  const tableChange = (pagination, filters, sort, extra) => {
    const { field, order } = sort;
    const { current, pageSize } = pagination;
    const obj = {
      ascend: 'asc',
      descend: 'desc',
    };

    if (order) {
      homeStore.SET_PROVIDERSPARAM({
        ...homeStore.params,
        page: current,
        page_size: pageSize,
        sort_by: field,
        order: obj[order],
      });
    } else {
      homeStore.SET_PROVIDERSPARAM({
        ...homeStore.params,
        sort_by: 'reputation_score',
        order: 'desc',
        page: current,
        page_size: pageSize,
      });
    }
    homeStore._getProviders(homeStore.params);
  };
  return (
    <div className={styles.list}>
      <Table
        columns={columns}
        dataSource={homeStore.list}
        rowKey={(record) => record.address}
        onChange={tableChange}
        scroll={{ x: '100%' }}
        pagination={{
          total: homeStore.total,
        }}
      />
    </div>
  );
};

export default observer(List);
