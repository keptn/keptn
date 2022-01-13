export interface EndSessionParameters {
  id_token_hint: string;
  state: string;
  post_logout_redirect_uri: string;
}

export interface EndSessionData extends EndSessionParameters {
  end_session_endpoint: string;
}
