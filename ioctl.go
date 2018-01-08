package evdev

import (
	"fmt"
	"syscall"
	"unsafe"
)

// ioctl wraps sending an ioctl syscall.
func ioctl(fd, name uintptr, data interface{}) error {
	var v uintptr

	switch dd := data.(type) {
	case unsafe.Pointer:
		v = uintptr(dd)

	case int:
		v = uintptr(dd)

	case uintptr:
		v = dd

	default:
		return fmt.Errorf("ioctl: data has invalid type %T", data)
	}

	_, _, errno := syscall.RawSyscall(syscall.SYS_IOCTL, fd, name, v)
	if errno != 0 {
		return errno
	}
	return nil
}

var (
	_EVIOCGVERSION    uintptr
	_EVIOCGID         uintptr
	_EVIOCGREP        uintptr
	_EVIOCSREP        uintptr
	_EVIOCGKEYCODE    uintptr
	_EVIOCGKEYCODE_V2 uintptr
	_EVIOCSKEYCODE    uintptr
	_EVIOCSKEYCODE_V2 uintptr
	_EVIOCSFF         uintptr
	_EVIOCRMFF        uintptr
	_EVIOCGEFFECTS    uintptr
	_EVIOCGRAB        uintptr
	_EVIOCSCLOCKID    uintptr

	// sizes
	sizeof_int    int
	sizeof_int2   int
	sizeof_id     int
	sizeof_keymap int
	sizeof_effect int
	sizeof_event  int
)

func init() {
	var i int32
	var id ID
	var keymap KeyMap
	var effect Effect
	var event Event

	sizeof_int = int(unsafe.Sizeof(i))
	sizeof_int2 = sizeof_int << 1
	sizeof_id = int(unsafe.Sizeof(id))
	sizeof_keymap = int(unsafe.Sizeof(keymap))
	sizeof_effect = int(unsafe.Sizeof(effect))
	sizeof_event = int(unsafe.Sizeof(event))

	_EVIOCGVERSION = _IOR('E', 0x01, sizeof_int)
	_EVIOCGID = _IOR('E', 0x02, sizeof_id)
	_EVIOCGREP = _IOR('E', 0x03, sizeof_int2)
	_EVIOCSREP = _IOW('E', 0x03, sizeof_int2)

	_EVIOCGKEYCODE = _IOR('E', 0x04, sizeof_int2)
	_EVIOCGKEYCODE_V2 = _IOR('E', 0x04, sizeof_keymap)
	_EVIOCSKEYCODE = _IOW('E', 0x04, sizeof_int2)
	_EVIOCSKEYCODE_V2 = _IOW('E', 0x04, sizeof_keymap)

	_EVIOCSFF = _IOC(_IOC_WRITE, 'E', 0x80, sizeof_effect)
	_EVIOCRMFF = _IOW('E', 0x81, sizeof_int)
	_EVIOCGEFFECTS = _IOR('E', 0x84, sizeof_int)
	_EVIOCGRAB = _IOW('E', 0x90, sizeof_int)
	_EVIOCSCLOCKID = _IOW('E', 0xa0, sizeof_int)
}

func _EVIOCGNAME(n int) uintptr {
	return _IOC(_IOC_READ, 'E', 0x06, n)
}

func _EVIOCGPHYS(n int) uintptr {
	return _IOC(_IOC_READ, 'E', 0x07, n)
}

func _EVIOCGUNIQ(n int) uintptr {
	return _IOC(_IOC_READ, 'E', 0x08, n)
}

func _EVIOCGPROP(n int) uintptr {
	return _IOC(_IOC_READ, 'E', 0x09, n)
}

func _EVIOCGMTSLOTS(n int) uintptr {
	return _IOC(_IOC_READ, 'E', 0x0a, n)
}

func _EVIOCGKEY(n int) uintptr {
	return _IOC(_IOC_READ, 'E', 0x18, n)
}

func _EVIOCGLED(n int) uintptr {
	return _IOC(_IOC_READ, 'E', 0x19, n)
}

func _EVIOCGSND(n int) uintptr {
	return _IOC(_IOC_READ, 'E', 0x1a, n)
}

func _EVIOCGSW(n int) uintptr {
	return _IOC(_IOC_READ, 'E', 0x1b, n)
}

func _EVIOCGBIT(ev, n int) uintptr {
	return _IOC(_IOC_READ, 'E', 0x20+ev, n)
}

func _EVIOCGABS(abs int) uintptr {
	var v Axis
	return _IOR('E', 0x40+abs, int(unsafe.Sizeof(v)))
}

func _EVIOCSABS(abs int) uintptr {
	var v Axis
	return _IOW('E', 0xc0+abs, int(unsafe.Sizeof(v)))
}

const (
	_IOC_NONE      = 0x0
	_IOC_WRITE     = 0x1
	_IOC_READ      = 0x2
	_IOC_NRBITS    = 8
	_IOC_TYPEBITS  = 8
	_IOC_SIZEBITS  = 14
	_IOC_DIRBITS   = 2
	_IOC_NRSHIFT   = 0
	_IOC_NRMASK    = (1 << _IOC_NRBITS) - 1
	_IOC_TYPEMASK  = (1 << _IOC_TYPEBITS) - 1
	_IOC_SIZEMASK  = (1 << _IOC_SIZEBITS) - 1
	_IOC_DIRMASK   = (1 << _IOC_DIRBITS) - 1
	_IOC_TYPESHIFT = _IOC_NRSHIFT + _IOC_NRBITS
	_IOC_SIZESHIFT = _IOC_TYPESHIFT + _IOC_TYPEBITS
	_IOC_DIRSHIFT  = _IOC_SIZESHIFT + _IOC_SIZEBITS
	_IOC_IN        = _IOC_WRITE << _IOC_DIRSHIFT
	_IOC_OUT       = _IOC_READ << _IOC_DIRSHIFT
	_IOC_INOUT     = (_IOC_WRITE | _IOC_READ) << _IOC_DIRSHIFT
	_IOCSIZE_MASK  = _IOC_SIZEMASK << _IOC_SIZESHIFT
)

func _IOC(dir, t, nr, size int) uintptr {
	return uintptr((dir << _IOC_DIRSHIFT) | (t << _IOC_TYPESHIFT) |
		(nr << _IOC_NRSHIFT) | (size << _IOC_SIZESHIFT))
}

func _IO(t, nr int) uintptr {
	return _IOC(_IOC_NONE, t, nr, 0)
}

func _IOR(t, nr, size int) uintptr {
	return _IOC(_IOC_READ, t, nr, size)
}

func _IOW(t, nr, size int) uintptr {
	return _IOC(_IOC_WRITE, t, nr, size)
}

func _IOWR(t, nr, size int) uintptr {
	return _IOC(_IOC_READ|_IOC_WRITE, t, nr, size)
}

func _IOC_DIR(nr int) uintptr {
	return uintptr(((nr) >> _IOC_DIRSHIFT) & _IOC_DIRMASK)
}

func _IOC_TYPE(nr int) uintptr {
	return uintptr(((nr) >> _IOC_TYPESHIFT) & _IOC_TYPEMASK)
}

func _IOC_NR(nr int) uintptr {
	return uintptr(((nr) >> _IOC_NRSHIFT) & _IOC_NRMASK)
}

func _IOC_SIZE(nr int) uintptr {
	return uintptr(((nr) >> _IOC_SIZESHIFT) & _IOC_SIZEMASK)
}
