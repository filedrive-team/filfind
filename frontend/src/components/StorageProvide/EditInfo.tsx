import storageProvideStore from '@/store/modules/storageProvide';
import { urlInterface } from '@/utils/interface';
import { Modal, Input, Form } from 'antd';
import { observer } from 'mobx-react';
import { useParams } from 'react-router';
enum isSPGroup {
  YES = 'Yes',
  NO = 'No',
}
const EditInfo = (props: any) => {
  const [form] = Form.useForm();
  const _url: urlInterface = useParams();
  const handleOk = async () => {
    const _form = form.getFieldsValue();
    form.validateFields(['is_member', 'parallel_deals']).then((res) => {
      const _params = {
        address: storageProvideStore.providerDetail.address,
        ..._form,
      };
      storageProvideStore._postProviderDetail(_params).then((data) => {
        storageProvideStore.SET_INFOVISIBLE(false);
        const address_id = _url.id;
        storageProvideStore._getSpServiceDetail({ address_id: address_id });
      });
    });
  };
  const handleCancel = () => {
    storageProvideStore.SET_INFOVISIBLE(false);
  };
  return (
    <Modal
      title={'EDIT ' + storageProvideStore.providerDetail.address}
      closable={false}
      visible={storageProvideStore.editInfoVisible}
      onOk={handleOk}
      onCancel={handleCancel}
      okText={'SAVE'}
      key="editInfo"
    >
      <Form
        layout="vertical"
        form={form}
        initialValues={storageProvideStore.providerDetail}
      >
        <Form.Item name={'available_deals'} label={'Available FIL+ deals'}>
          <Input placeholder="e.g., 10TiB/D" />
        </Form.Item>

        <Form.Item name={'bandwidth'} label={'Data transfer bandwidth'}>
          <Input placeholder="e.g., 300Mbps" />
        </Form.Item>

        <div className="inputRow two">
          <Form.Item name={'sealing_speed'} label={'Sealing speed'}>
            <Input placeholder="e.g., 10TiB/D" />
          </Form.Item>

          <Form.Item name={'renewable_energy'} label={'Renewable energy'}>
            <Input placeholder="e.g., 1MWh" />
          </Form.Item>
        </div>

        <Form.Item
          name={'parallel_deals'}
          label={'Number of parallel deals'}
          rules={[
            {
              validator: (_, value) => {
                const reg = /^[0-9]+$/;
                if (value) {
                  const valid = reg.test(value);
                  if (valid) {
                    return Promise.resolve();
                  } else {
                    return Promise.reject(
                      new Error('Please enter a positive integer!'),
                    );
                  }
                } else {
                  return Promise.resolve();
                }
              },
            },
          ]}
        >
          <Input placeholder="e.g., 10" />
        </Form.Item>

        <Form.Item name={'certification'} label={'Certification'}>
          <Input placeholder="e.g., PCI Compliance" />
        </Form.Item>

        <Form.Item
          name={'is_member'}
          label={'Member of Enterprise SP Group'}
          rules={[
            {
              validator: (_, value) => {
                if (value) {
                  const _list = [isSPGroup.YES, isSPGroup.NO];
                  const valid = _list.includes(value);
                  if (valid) {
                    return Promise.resolve();
                  } else {
                    return Promise.reject(new Error('Please enter Yes or No!'));
                  }
                } else {
                  return Promise.resolve();
                }
              },
            },
          ]}
        >
          <Input placeholder="e.g., Yes or No" />
        </Form.Item>

        <Form.Item name={'experience'} label={'Onboarding Products Experience'}>
          <Input placeholder="e.g., Textile, Estuary, Slingshot, etc." />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default observer(EditInfo);
