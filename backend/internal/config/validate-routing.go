package config

import "fmt"

func (s *ServiceConfig) Validate() []error {
	var errs []error

	if s == nil {
		errs = append(errs, fmt.Errorf("service configuration is nil"))
		return errs
	}

	if s.ServiceName == "" {
		errs = append(errs, fmt.Errorf("service_name must not be empty"))
	}

	if s.URL == "" {
		errs = append(errs, fmt.Errorf("url must not be empty for service: %s", s.ServiceName))
	}

	if len(s.SupportedVersions) == 0 {
		errs = append(errs, fmt.Errorf("supported_versions must contain at least one version for service: %s", s.ServiceName))
	} else {
		for _, version := range s.SupportedVersions {
			if version == "" {
				errs = append(errs, fmt.Errorf("supported_versions contains an empty version string for service: %s", s.ServiceName))
			}
		}
	}

	return errs
}

func (r *RouteConfig) Validate() []error {
	var errs []error

	if r == nil {
		errs = append(errs, fmt.Errorf("route configuration is nil"))
		return errs
	}

	if r.Path == "" {
		errs = append(errs, fmt.Errorf("path must not be empty"))
	}

	if len(r.Methods) > 0 {
		for i, method := range r.Methods {
			if method == "" {
				errs = append(errs, fmt.Errorf("path %s, method %d is empty string", r.Path, i+1))
			}
		}
	}

	if len(r.Services) == 0 {
		errs = append(errs, fmt.Errorf("no services defined for path: %s", r.Path))
	} else {
		for i, service := range r.Services {
			serviceErrs := service.Validate()
			for _, err := range serviceErrs {
				errs = append(errs, fmt.Errorf("path %s, service %d: %v", r.Path, i+1, err))
			}
		}
	}

	return errs
}
