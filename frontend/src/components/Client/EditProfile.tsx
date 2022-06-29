import { Modal, Form } from 'antd';
import { Input } from 'antd';
import clientStore from '@/store/modules/client';

import { observer } from 'mobx-react';
import Upload from '@/components/Upload';
import { Cookies } from '@/data/cookies';
import { urlInterface } from '@/utils/interface';
import { useParams } from 'react-router';
const EditProfile = (props: any) => {
  const _url: urlInterface = useParams();
  const [form] = Form.useForm();
  const handleOk = () => {
    form.validateFields().then((data) => {
      const _form = form.getFieldsValue();
      const _params = {
        ..._form,
        avatar: clientStore.clientProfile.avatar,
        logo: clientStore.clientProfile.logo,
      };
      clientStore._postUserProfile(_params).then((res) => {
        clientStore.SET_PROFILEVISIBLE(false);
        const _userInfo = Cookies.getUserInfo();
        if (_userInfo) {
          Cookies.setUserInfo({
            ..._userInfo,
            name: clientStore.clientProfile.name,
            email: clientStore.clientProfile.contact_email,
            location: clientStore.clientProfile.location,
            avatar: clientStore.clientProfile.avatar,
            logo: clientStore.clientProfile.logo,
          });
        }
        const address_id = _url.id;
        clientStore._getClientProfile({ address_id: address_id });
        clientStore.SET_PROFILEVISIBLE(false);
      });
    });
  };
  const handleCancel = () => {
    clientStore.SET_PROFILEVISIBLE(false);
  };
  const avatarChange = (value) => {
    clientStore.SET_CLIENTPROFILE({
      ...clientStore.clientProfile,
      avatar: value,
    });
  };

  const logoChange = (value) => {
    clientStore.SET_CLIENTPROFILE({
      ...clientStore.clientProfile,
      logo: value,
    });
  };
  return (
    <Modal
      title="Edit Profile"
      closable={false}
      visible={clientStore.editProfileVisible}
      onOk={handleOk}
      onCancel={handleCancel}
      okText={'SAVE'}
      key="editProfile"
    >
      <Form
        form={form}
        layout="vertical"
        initialValues={clientStore.clientProfile}
      >
        <div className="inputRow upload">
          <Upload
            hash={clientStore.clientProfile.avatar}
            label={'avatar'}
            valueChange={avatarChange}
          ></Upload>
          <Upload
            hash={clientStore.clientProfile.logo}
            label={'logo'}
            valueChange={logoChange}
          ></Upload>
        </div>
        <div className="inputRow">
          <Form.Item
            name={'name'}
            label={'Name'}
            rules={[{ required: true, message: 'Please Enter Name!' }]}
          >
            <Input placeholder="Input your name" />
          </Form.Item>
          <Form.Item name={'location'} label={'Location'}>
            <Input placeholder="Input your location" />
          </Form.Item>
        </div>
        <Form.Item name={'description'} label={'Brief Introduction'}>
          <Input placeholder="Input your Bio" />
        </Form.Item>
        <Form.Item
          name={'contact_email'}
          label={'Email'}
          rules={[{ type: 'email' }]}
        >
          <Input placeholder="Enter your e-mail" />
        </Form.Item>
        <Form.Item name={'slack'} label={'Slack'}>
          <Input placeholder="Input your Slack handle" />
        </Form.Item>
        <Form.Item name={'twitter'} label={'Twitter'}>
          <Input placeholder="Input your Twitter handle" />
        </Form.Item>
        <Form.Item name={'github'} label={'GitHub'}>
          <Input placeholder="Input your GitHub handle" />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default observer(EditProfile);
