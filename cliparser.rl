package main

%%{
    machine cliparser;
    write data;
}%%

type Cli struct {
    status int
    body []byte
}

func Cliparser(data []byte) (cli *Cli){
    cs, p, pe := 0, 0, len(data)
    cli = new(Cli)
    bodylength, bodypos := 0, 0
    %%{
        action status {cli.status = cli.status*10+(int(fc)-'0')}
        action bodylength {bodylength = bodylength*10+(int(fc)-'0')}
        action makebody {cli.body = make([]byte,bodylength)}
        action body {cli.body[bodypos]=fc; bodypos++; if bodypos == bodylength {fbreak;}}
        main := digit{,3}@status " " digit+ @bodylength %makebody space* "\n" (any*)@body;
        write init;
        write exec;
    }%%

    return cli
}

