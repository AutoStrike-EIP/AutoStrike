package com.autostrike.backend.campaign;

import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.ResponseStatus;

/** La campagne n'a pas pu être transmise au Task Manager (broker injoignable). */
@ResponseStatus(HttpStatus.SERVICE_UNAVAILABLE)
public class CampaignLaunchException extends RuntimeException {

  public CampaignLaunchException(Throwable cause) {
    super("Message broker unavailable, the scenario could not be launched", cause);
  }
}
