gir
===

This repository provides GIR-based autogenerated Go bindings for GObject and GObject-related libraries, such as Gtk4. These bindings are significantly lower level than most, preferring to allow the user to perform manual memory management rather than attempting to hook the GObject system up to Go's garbage collector. If you are in doubt about which bindings to use, these are probably not the ones that you want.

_Disclaimer: The above claim is technically inaccurate as the implementation is not to the point of usability yet._
