import styles from './List.module.scss';
import { Table } from 'antd';
import homeStore from '@/store/modules/home';
import { observer } from 'mobx-react';
import { bytesToSize } from '@/utils';
import clientListStore from '@/store/modules/clientList';

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
      title: 'Client ID',
      dataIndex: 'address_id',
      key: 'address_id',
      render: (text, record) => {
        const { address_id } = record;
        return (
          <div onClick={() => skipToClient(address_id)} className="listLink">
            {text}
          </div>
        );
      },
    },
    {
      title: 'Client Name',
      dataIndex: 'name',
      key: 'name',
      render: (text) => {
        return <div>{text === '' ? '-' : text}</div>;
      },
    },
    {
      title: 'Location',
      dataIndex: 'location',
      key: 'location',
      render: (text) => {
        return <div>{text === '' ? '-' : text}</div>;
      },
    },
    {
      title: 'Total Storage Deals',
      dataIndex: 'storage_deals',
      key: 'storage_deals',
      sorter: true,
      render: (text, record) => {
        return <div>{text}</div>;
      },
    },
    {
      title: 'Storage Capacity',
      dataIndex: 'storage_capacity',
      key: 'storage_capacity',
      sorter: true,
      render: (text, record) => {
        const _format = bytesToSize(text);
        return <div>{_format}</div>;
      },
    },
    {
      title: 'Allocated DataCap',
      dataIndex: 'used_data_cap',
      key: 'used_data_cap',
      sorter: true,
      render: (text, record) => {
        const _format = bytesToSize(text);
        return <div>{_format}</div>;
      },
    },
    {
      title: 'Available DataCap',
      dataIndex: 'data_cap',
      key: 'data_cap',
      sorter: true,
      render: (text, record) => {
        const _format = bytesToSize(text);
        return <div>{_format}</div>;
      },
    },
    {
      title: 'Monthly Storage Capacity',
      dataIndex: 'monthly_storage',
      key: 'monthly_storage',
      render: (text) => {
        return <div>{text === '' ? '-' : text}</div>;
      },
    },
    {
      title: 'Data Transfer Bandwidth',
      dataIndex: 'bandwidth',
      key: 'bandwidth',
      render: (text) => {
        return <div>{text === '' ? '-' : text}</div>;
      },
    },
  ];
  const skipToClient = (id) => {
    window.open(`/client/${id}`);
  };
  const tableChange = (pagination, filters, sort, extra) => {
    const { field, order } = sort;
    const obj = {
      ascend: 'asc',
      descend: 'desc',
    };
    if (order) {
      clientListStore.SET_CLIENTSPARAM({
        ...clientListStore.params,
        sort_by: field,
        order: obj[order],
      });
    } else {
      clientListStore.SET_CLIENTSPARAM({
        ...clientListStore.params,
        sort_by: 'data_cap',
        order: 'asc',
      });
    }
    clientListStore._getClients(clientListStore.params);
  };

  return (
    <div className={styles.list}>
      <div className={styles.listContent}>
        <Table
          columns={columns}
          dataSource={clientListStore.list}
          rowKey={(record) => record.address_id}
          onChange={tableChange}
          pagination={{
            total: clientListStore.total,
            onChange: (page, pageSize) => {
              clientListStore.SET_CLIENTSPARAM({
                ...homeStore.params,
                page: page,
                page_size: pageSize,
              });
              clientListStore._getClients(clientListStore.params);
            },
          }}
        />
      </div>
    </div>
  );
};

export default observer(List);
