import HttpErrorContent from '@/components/HttpErrorContent';

const Index = ({ errCode = '50X', errMsg = '' }) => {
  return <HttpErrorContent httpCode={errCode} errMsg={errMsg} />;
};

export default Index;
