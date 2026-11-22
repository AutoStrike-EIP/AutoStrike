package com.autostrike.backend.campaign;

import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.ResponseStatus;

@ResponseStatus(HttpStatus.BAD_REQUEST)
public class UnknownTemplateException extends RuntimeException {

  public UnknownTemplateException(String templateId) {
    super("Unknown scenario template: " + templateId);
  }
}
