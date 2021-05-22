Easy I2P Zero Bundling Tool for Go
==================================

This accomplishes the bundling and managment of an I2P Zero router from within a Go application. The recommended use of
this library is as a way of providing an embedded *option* in an out-of-tree application, an out-of-tree application
which would otherwise prefer an existing I2P Router. This means that the standard way to use it should be:

1. Check for a running I2P Router. If one is found, skip to step *5.*

2. Check for a stopped I2P Router installed from an official package.

3. If a stopped package-installed I2P router is available from an official package, start it. Skip step *4.*

4. If no other I2P Router is available, start the embedded I2P Zero Router.

5. Start the external application.

Doing things in this way allows us to conserve resources by not running redundant I2P Routers on the same computer,
while also allowing the use of an embedded I2P router to auto-configure standalone applications on computer even when
I2P is not present.

Use Scenarios
-------------

### **Scenario A:** I2P Router is installed on host PC *Prior to* the first run of the out-of-tree application.

In this scenario, the I2P Router in use is the package-installed router, and the embedded one is left alone.

### **Scenario B:** out-of-tree application is installed on host PC *alone*, with no other router available

In this scenario, the I2P Router in use is the embedded router.

### **Scenario C:** out-of-tree application is installed on host PC *prior to* a system-wide I2P router which becomes preferred.

In this scenario, the keys used for identifying the SAM tunnels are managed by the application and thus migrate with
the application from the embedded router to the package-installed router.

Example Usage:
--------------

        package main

        import (
          "log"
        )

        import (
          "github.com/eyedeekay/zerobundle"
        )

        func main() {
          if err := zerobundle.ZeroMain(); err != nil {
            log.Println(err)
          }
        }

