mongo -- "mongoadmin" <<EOF
var user = 'mongoadmin';
var passwd = 'mongoadmin';
var admin = db.getSiblingDB('admin');
admin.auth(user, passwd);
db.createUser({user: user, pwd: passwd, roles: ["readWrite"]});
EOF