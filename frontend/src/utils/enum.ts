enum SignType {
  signIn = 'signIn',
  signUp = 'signUp',
}

enum SPSStatusType {
  All = 'all',
  Registered = 'registered',
  Autofilled = 'autofilled',
}

enum RegionType {
  All = 'all',
  Africa = 'Africa',
  Asia = 'Asia',
  Europe = 'Europe',
  NorthAmerica = 'North America',
  Oceania = 'Oceania',
  SouthAmerica = 'South America',
}

enum ValidateStatus {
  none = '',
  success = '',
  warning = 'warning',
  error = 'error',
  validating = 'validating',
}

enum UserType {
  sp = 'sp_owner',
  client = 'data_client',
}

enum envType {
  production = 'production',
  development = 'development',
  testnet = 'testnet',
}
export {
  SignType,
  SPSStatusType,
  RegionType,
  ValidateStatus,
  UserType,
  envType,
};
