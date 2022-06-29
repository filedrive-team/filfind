import { Button, Result } from 'antd';
import { useHistory } from 'react-router';

const NoFoundPage = () => {
  const history = useHistory();
  return (
    <Result
      status="404"
      subTitle="The page you visited does not exist."
      extra={
        <Button type="primary" onClick={() => history.push('/')}>
          Back Home
        </Button>
      }
    />
  );
};

export default NoFoundPage;
