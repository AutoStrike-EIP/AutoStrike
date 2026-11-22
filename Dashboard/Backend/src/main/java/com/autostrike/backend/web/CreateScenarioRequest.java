package com.autostrike.backend.web;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;

public record CreateScenarioRequest(@NotBlank @Size(max = 255) String name) {}
