// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2016 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package builtin

import (
	"github.com/ubuntu-core/snappy/interfaces"
)

// http://bazaar.launchpad.net/~ubuntu-security/ubuntu-core-security/trunk/view/head:/data/apparmor/policygroups/ubuntu-core/16.04/unity7
const unity7ConnectedPlugAppArmor = `
# Description: Can access Unity7. Restricted because Unity 7 runs on X and
# requires access to various DBus services and this enviroment does not prevent
# eavesdropping or apps interfering with one another.
# Usage: reserved

#include <abstractions/dbus-strict>
#include <abstractions/dbus-session-strict>
#include <abstractions/X>

#include <abstractions/fonts>
/var/cache/fontconfig/   r,
/var/cache/fontconfig/** mr,

# subset of gnome abstraction
/etc/gnome/defaults.list r,
/usr/share/gnome/applications/             r,
/usr/share/applications/mimeinfo.cache     r,

/etc/gtk-*/*                               r,
/usr/lib{,32,64}/gtk-*/**                  mr,
/usr/lib{,32,64}/gdk-pixbuf-*/**           mr,
/usr/lib/@{multiarch}/gtk-*/**             mr,
/usr/lib/@{multiarch}/gdk-pixbuf-*/**      mr,

/etc/pango/*                               r,
/usr/lib{,32,64}/pango/**                  mr,
/usr/lib/@{multiarch}/pango/**             mr,

/usr/share/icons/                          r,
/usr/share/icons/**                        r,
/usr/share/icons/*/index.theme             rk,
/usr/share/pixmaps/                        r,
/usr/share/pixmaps/**                      r,
/usr/share/unity/icons/**                  r,
/usr/share/thumbnailer/icons/**            r,
/usr/share/themes/**                       r,

#owner @{HOME}/.themes/                r,
#owner @{HOME}/.themes/**              r,


# subset of ibus abstraction
/usr/lib/@{multiarch}/gtk-2.0/[0-9]*/immodules/im-ibus.so mr,
owner @{HOME}/.config/ibus/      r,
owner @{HOME}/.config/ibus/bus/  r,
owner @{HOME}/.config/ibus/bus/* r,

# subset of freedesktop.org
/usr/share/mime/**                   r,
owner @{HOME}/.local/share/mime/**   r,
owner @{HOME}/.config/user-dirs.dirs r,

# accessibility
#include <abstractions/dbus-accessibility-strict>
dbus (send)
    bus=session
    path=/org/a11y/bus
    interface=org.a11y.Bus
    member=GetAddress
    peer=(label=unconfined),

# unfortunate, but org.a11y.atspi is not designed for separation
dbus (receive, send)
    bus=accessibility
    path=/org/a11y/atspi/**
    peer=(label=unconfined),

# org.freedesktop.Accounts
dbus (send)
    bus=system
    path=/org/freedesktop/Accounts
    interface=org.freedesktop.DBus.Introspectable
    member=Introspect
    peer=(label=unconfined),

dbus (send)
    bus=system
    path=/org/freedesktop/Accounts
    interface=org.freedesktop.Accounts
    member=FindUserById
    peer=(label=unconfined),

# Get() is an information leak
# TODO: verify what it is leaking
dbus (receive, send)
    bus=system
    path=/org/freedesktop/Accounts/User[0-9]*
    interface=org.freedesktop.DBus.Properties
    member={Get,PropertiesChanged}
    peer=(label=unconfined),

# gmenu
# Note: the gmenu DBus api was not designed for application isolation and apps
# may specify anything as their 'path'. For example, these work in the many
# cases:
# - /org/gtk/Application/anonymous{,/**}
# - /com/canonical/unity/gtk/window/[0-9]*
# but libreoffice does:
# - /org/libreoffice{,/**}
# As such, cannot mediate by DBus path so we'll be as strict as we can in the
# other mediated parts
dbus (send)
    bus=session
    interface=org.gtk.Actions
    member=Changed
    peer=(name=org.freedesktop.DBus, label=unconfined),

dbus (receive)
    bus=session
    interface=org.gtk.Actions
    member={Activate,DescribeAll,SetState}
    peer=(label=unconfined),

dbus (receive)
    bus=session
    interface=org.gtk.Menus
    member={Start,End}
    peer=(label=unconfined),

dbus (send)
    bus=session
    interface=org.gtk.Menus
    member=Changed
    peer=(name=org.freedesktop.DBus, label=unconfined),

# url helper
dbus (send)
    bus=session
    interface=com.canonical.UrlLauncher.XdgOpen
    peer=(label=unconfined),

# dbusmenu
dbus (send)
    bus=session
    path=/MenuBar
    interface=com.canonical.dbusmenu
    member="{LayoutUpdated,ItemsPropertiesUpdated}"
    peer=(name=org.freedesktop.DBus, label=unconfined),

dbus (receive)
    bus=session
    path=/MenuBar
    interface="{com.canonical.dbusmenu,org.freedesktop.DBus.Properties}"
    member=Get*
    peer=(label=unconfined),

dbus (receive)
    bus=session
    path=/MenuBar
    interface=com.canonical.dbusmenu
    member="{AboutTo*,Event*}"
    peer=(label=unconfined),

# notifications
dbus (send)
    bus=session
    path=/StatusNotifierWatcher
    interface=org.freedesktop.DBus.Introspectable
    member=Introspect
    peer=(name=org.kde.StatusNotifierWatcher, label=unconfined),

dbus (send)
    bus=session
    path=/org/freedesktop/DBus
    interface=org.freedesktop.DBus
    member="{GetConnectionUnixProcessID,RequestName,ReleaseName}"
    peer=(name=org.freedesktop.DBus, label=unconfined),

dbus (bind)
    bus=session
    name=org.kde.StatusNotifierItem-[0-9]*,

dbus (send)
    bus=session
    path=/StatusNotifierWatcher
    interface=org.freedesktop.DBus.Properties
    member=Get
    peer=(name=org.kde.StatusNotifierWatcher, label=unconfined),

dbus (send)
    bus=session
    path=/StatusNotifierWatcher
    interface=org.kde.StatusNotifierWatcher
    member=RegisterStatusNotifierItem
    peer=(name=org.kde.StatusNotifierWatcher, label=unconfined),

dbus (send)
    bus=session
    path=/StatusNotifierItem
    interface=org.kde.StatusNotifierItem
    member="New{AttentionIcon,Icon,OverlayIcon,Status,Title,ToolTip}"
    peer=(name=org.freedesktop.DBus, label=unconfined),

dbus (receive)
    bus=session
    path=/StatusNotifierItem
    interface=org.freedesktop.DBus.Properties
    member=Get*
    peer=(label=unconfined),

dbus (send)
    bus=session
    path=/org/freedesktop/Notifications
    interface=org.freedesktop.Notifications
    member="{GetCapabilities,GetServerInformation,Notify}"
    peer=(name=org.freedesktop.Notifications, label=unconfined),

dbus (receive)
    bus=session
    path=/org/freedesktop/Notifications
    interface=org.freedesktop.Notifications
    member=NotificationClosed
    peer=(label=unconfined),
>>>>>>> upstream/master

# unity launcher
dbus (send)
    bus=session
    path=/com/canonical/unity/launcherentry/[0-9]*
    interface=com.canonical.Unity.LauncherEntry
    member=Update
    peer=(name=org.freedesktop.DBus, label=unconfined),

dbus (send)
    bus=session
    path=/com/canonical/unity/launcherentry/[0-9]*
    interface=com.canonical.dbusmenu
    member="{LayoutUpdated,ItemsPropertiesUpdated}"
    peer=(name=org.freedesktop.DBus, label=unconfined),

dbus (receive)
    bus=session
    path=/com/canonical/unity/launcherentry/[0-9]*
    interface="{com.canonical.dbusmenu,org.freedesktop.DBus.Properties}"
    member=Get*
    peer=(label=unconfined),

# This rule is meant to be covered by abstractions/dbus-session-strict but
# the unity launcher code has a typo that uses /org/freedesktop/dbus as the
# path instead of /org/freedesktop/DBus, so we need to all it here.
dbus (send)
    bus=session
    path=/org/freedesktop/dbus
    interface=org.freedesktop.DBus
    member=NameHasOwner
    peer=(name=org.freedesktop.DBus, label=unconfined),

# appmenu
dbus (send)
    bus=session
    path=/org/freedesktop/DBus
    interface=org.freedesktop.DBus
    member=ListNames
    peer=(name=org.freedesktop.DBus, label=unconfined),

dbus (send)
    bus=session
    path=/com/canonical/AppMenu/Registrar
    interface=com.canonical.AppMenu.Registrar
    member="{RegisterWindow,UnregisterWindow}"
    peer=(label=unconfined),

dbus (send)
    bus=session
    path=/com/canonical/AppMenu/Registrar
    interface=com.canonical.dbusmenu
    member=UnregisterWindow
    peer=(label=unconfined),

dbus (receive)
    bus=session
    path=/com/canonical/menu/[0-9]*
    interface="{org.freedesktop.DBus.Properties,com.canonical.dbusmenu}"
    member="{GetAll,GetLayout}"
    peer=(label=unconfined),


# Lttng tracing is very noisy and should not be allowed by confined apps. Can
# safely deny. LP: #1260491
deny /{,var/}run/shm/lttng-ust-* r,
`

// http://bazaar.launchpad.net/~ubuntu-security/ubuntu-core-security/trunk/view/head:/data/seccomp/policygroups/ubuntu-core/16.04/unity7
const unity7ConnectedPlugSecComp = `
# Description: Can access Unity7. Restricted because Unity 7 runs on X and
# requires access to various DBus services and this enviroment does not prevent
# eavesdropping or apps interfering with one another.

# X
getpeername
recvfrom
recvmsg
shutdown
getsockopt

# dbus
connect
getsockname
recvmsg
send
sendto
sendmsg
socket
`

// NewUnity7Interface returns a new "unity7" interface.
func NewUnity7Interface() interfaces.Interface {
	return &commonInterface{
		name: "unity7",
		connectedPlugAppArmor: unity7ConnectedPlugAppArmor,
		connectedPlugSecComp:  unity7ConnectedPlugSecComp,
		reservedForOS:         true,
		autoConnect:           true,
	}
}
