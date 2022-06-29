import { Modal, Form, Input } from 'antd';
import clientStore from '@/store/modules/client';
import { observer } from 'mobx-react';
import { urlInterface } from '@/utils/interface';
import { useParams } from 'react-router';
const EditInfo = (props: any) => {
  const _params: urlInterface = useParams();
  const [form] = Form.useForm();
  const handleOk = () => {
    const _form = form.getFieldsValue();
    clientStore._postClientDetail(_form).then((res) => {
      clientStore.SET_INFOVISIBLE(false);
      const address_id = _params.id;
      clientStore._getClientDetail({ address_id: address_id });
    });
  };
  const handleCancel = () => {
    clientStore.SET_INFOVISIBLE(false);
  };
  return (
    <Modal
      title={'Edit ' + clientStore.clientDetail.address_id}
      closable={false}
      visible={clientStore.editInfoVisible}
      onOk={handleOk}
      onCancel={handleCancel}
      okText={'SAVE'}
      key="editInfo"
    >
      <Form
        form={form}
        layout="vertical"
        initialValues={clientStore.clientDetail}
      >
        <Form.Item
          name={'monthly_storage'}
          label={'Monthly Storage requirement'}
        >
          <Input placeholder="e.g., 10TiB" />
        </Form.Item>
        <Form.Item name={'bandwidth'} label={'Data transfer bandwidth'}>
          <Input placeholder="e.g., 300Mbps" />
        </Form.Item>

        <Form.Item name={'use_case'} label={'Data use case'}>
          <Input placeholder="e.g., Entertainment / Media / Science" />
        </Form.Item>
        <Form.Item
          name={'service_requirement'}
          label={'Storage Service Requirement'}
        >
          <Input placeholder="er your req" />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default observer(EditInfo);
