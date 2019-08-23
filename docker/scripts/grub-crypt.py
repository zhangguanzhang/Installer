#! /usr/bin/python

'''Generate encrypted passwords for GRUB.'''

import crypt
import getopt
import getpass
import sys

def usage():
    '''Output usage message to stderr and exit.'''
    print >> sys.stderr, 'Usage: grub-crypt [OPTION]...'
    print >> sys.stderr, 'Try `$progname --help\' for more information.'
    sys.exit(1)

def gen_salt():
    '''Generate a random salt.'''
    ret = ''
    with open('/dev/urandom', 'rb') as urandom:
        while True:
            byte = urandom.read(1)
            if byte in ('ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz'
                        './0123456789'):
                ret += byte
                if len(ret) == 16:
                    break
    return ret

def main():
    '''Top level.'''
    crypt_type = '$6$' # SHA-256
    try:
        opts, args = getopt.getopt(sys.argv[1:], 'hv',
                                   ('help', 'version', 'md5', 'sha-256',
                                    'sha-512'))
    except getopt.GetoptError, err:
        print >> sys.stderr, str(err)
        usage()
    if args:
        print >> sys.stderr, 'Unexpected argument `%s\'' % (args[0],)
        usage()
    for (opt, _) in opts:
        if opt in ('-h', '--help'):
            print (
'''Usage: grub-crypt [OPTION]...
Encrypt a password.
  -h, --help              Print this message and exit
  -v, --version           Print the version information and exit
  --md5                   Use MD5 to encrypt the password
  --sha-256               Use SHA-256 to encrypt the password
  --sha-512               Use SHA-512 to encrypt the password (default)
Report bugs to <bug-grub@gnu.org>.
EOF''')
            sys.exit(0)
        elif opt in ('-v', '--version'):
            print 'grub-crypt (GNU GRUB 0.97)'
            sys.exit(0)
        elif opt == '--md5':
            crypt_type = '$1$'
        elif opt == '--sha-256':
            crypt_type = '$5$'
        elif opt == '--sha-512':
            crypt_type = '$6$'
        else:
            assert False, 'Unhandled option'
    # Modified to accept non-TTY input, though this is not recommended
    # for production use.
    if sys.stdin.isatty():
        password = getpass.getpass('Password: ')
        password2 = getpass.getpass('Retype password: ')
    else:
        print 'Password: '
        password = sys.stdin.readline().rstrip()
        print 'Retype password: '
        password2 = sys.stdin.readline().rstrip()

    if not password:
        print >> sys.stderr, 'Empty password is not permitted.'
        sys.exit(1)
    if password != password2:
        print >> sys.stderr, 'Sorry, passwords do not match.'
        sys.exit(1)
    salt = crypt_type + gen_salt()
    print crypt.crypt(password, salt)

if __name__ == '__main__':
    main()
