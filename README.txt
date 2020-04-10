
package pager // import "github.com/jackdoe/go-pager"


FUNCTIONS

func Pager(try ...string) (io.Writer, func())
    Create new pager executing a command based on $PAGER env var or array of
    executables Example:

        p, close := pager.Pager("less", "more","cat")
        defer close()

        p.Write("hello world")

    Will try to find $PAGER,less,more,cat in path, first one it finds it will
    pipe the output written to the returned Writer, if nothing is found it will
    return os.Stdout, if $PAGER=NOPAGER it will return os.Stdout, if $PAGER is
    specified but cant be found in path it will panic

