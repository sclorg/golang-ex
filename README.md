## Hello, OpenShift! ##

A sample app that is built using the [S2I golang builder](https://github.com/sclorg/golang-container).

This example will serve an HTTP response of "Hello OpenShift!" written in Golang. It is also
intended to be used with an evolving [Golang Source-to-Image builder image](https://github.com/sclorg/golang-container).

Once the image-stream is loaded into OpenShift, you can easily deploy it by doing:

    $ oc new-app golang~https://github.com/sclorg/golang-ex.git

Note: this is reused [example hello_openshift from OpenShift Origin](https://github.com/openshift/origin), separating it out will allow only the need to clone this example repo instead of all of the origin one.

The response message can be set by using the RESPONSE environment
variable.  You will need to edit the pod definition and add an
environment variable to the container definition and run the new pod.

Then you can re-create the pod as with the first example, get the new IP
address, and then curl will show your new message:

    $ curl 10.1.0.2:8888
     Hello World!

To test from external network, you need to create router. Please refer to [Running the router](https://github.com/openshift/origin/blob/master/docs/routing.md)

The container doesn't expose any ports, this will require you to expose it manually.
For example:

    $ oc expose dc/golang-ex --port=8888

and lastly if you want to expose a route, by doing:

    $ oc expose service/golang-ex


