set ver=Vers#v%date:~0,4%%date:~5,2%%date:~8,2%%time:~0,2%%time:~3,2%%time:~6,2%
.\windows\proto_gen_ext.exe --startid=100 --proto=pb --in=.\..\example\ --out=.\..\output\ --type=csharp --idrule=hash --idtype=uint16 --version=%ver%