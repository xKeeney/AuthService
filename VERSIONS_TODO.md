# Future features (road to v1.0.0 RELEASE)
- [ ] WebUI (client)
- [ ] Roles + permissions (admin / user / support / service) | WebUI first
- [ ] Service administration | WebUI first.
- [ ] Maintenance mode (middleware that dont allow new connections in maintenance mode)
- [ ] Update mode (associations) | associations after roles + permissions
- [ ] Create admin user in first run (if no admin)

# v0.3.0 beta "AuditLog"
- [ ] Add audit log (database + logic)
- [ ] Register audit log
- [ ] Login audit log
- [ ] Logout audit log
- [ ] Refresh audit log
- [ ] Revoke audit log
- [ ] Security event if refresh with revoked token

# v0.2.0 beta "UpdateMode"
- [ ] Add update mode
- [ ] Add migrations in update mode

# v0.1.0 beta "MVP"
- [x] add **Register** endpoint + logic | __26.02.2026__
- [x] Generate Keys for JWT mode | __28.02.2026__
- [ ] add JWT tokens logic (RS256 algorithm, refresh tokens can be used only one time)
    - [x] Access tokens | __28.02.2026__
    - [ ] Refresh tokens
- [x] add **Login** endpoint + logic __28.02.2026__
    - [x] Return access token | __28.02.2026__
    - [x] Add refresh token in cookie __28.02.2026__
- [ ] add **Refresh** endpoint + logic
- [ ] add **Logout** endpoint + logic
