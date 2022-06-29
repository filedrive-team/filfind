import { Form, Modal } from 'antd';
import { Input } from 'antd';
import { observer } from 'mobx-react';
import Upload from '@/components/Upload';
import storageProvideStore from '@/store/modules/storageProvide';
import { Cookies } from '@/data/cookies';
import { useParams } from 'react-router';
import { urlInterface } from '@/utils/interface';

const EditProfile = (props: any) => {
  const _url: urlInterface = useParams();
  const [form] = Form.useForm();
  const handleOk = () => {
    form.validateFields().then((data) => {
      const _form = form.getFieldsValue();
      const _params = {
        ..._form,
        avatar: storageProvideStore.userProfile.avatar,
        logo: storageProvideStore.userProfile.logo,
      };
      storageProvideStore._postUserProfile(_params).then((res) => {
        storageProvideStore.SET_PROFILEVISIBLE(false);
        const _userInfo = Cookies.getUserInfo();
        if (_userInfo) {
          Cookies.setUserInfo({
            ..._userInfo,
            name: storageProvideStore.userProfile.name,
            email: storageProvideStore.userProfile.contact_email,
            location: storageProvideStore.userProfile.location,
            avatar: storageProvideStore.userProfile.avatar,
            logo: storageProvideStore.userProfile.logo,
          });
        }
        const address_id = _url.id;
        storageProvideStore._getSpOwnerProfile({ address_id: address_id });
      });
    });
  };
  const handleCancel = () => {
    storageProvideStore.SET_PROFILEVISIBLE(false);
  };

  const avatarChange = (value) => {
    storageProvideStore.SET_USERPROFILE({
      ...storageProvideStore.userProfile,
      avatar: value,
    });
  };

  const logoChange = (value) => {
    storageProvideStore.SET_USERPROFILE({
      ...storageProvideStore.userProfile,
      logo: value,
    });
  };

  return (
    <Modal
      title="Edit Profile"
      closable={false}
      visible={storageProvideStore.editProfileVisible}
      onOk={handleOk}
      onCancel={handleCancel}
      okText={'SAVE'}
      key="editProfile"
    >
      <Form
        form={form}
        layout="vertical"
        initialValues={storageProvideStore.userProfile}
      >
        <div className="inputRow upload">
          <Upload
            hash={storageProvideStore.userProfile.avatar}
            label={'avatar'}
            valueChange={avatarChange}
          ></Upload>
          <Upload
            hash={storageProvideStore.userProfile.logo}
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
          rules={[
            { type: 'email', message: 'Please enter the correct email.' },
          ]}
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
