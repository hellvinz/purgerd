package client

var _cliparser_actions [] int8  = [] int8  { 0, 1, 0, 1, 1, 1, 2, 1, 3, 0  }
var _cliparser_key_offsets [] int8  = [] int8  { 0, 0, 3, 5, 11, 15, 18, 21, 22, 0  }
var _cliparser_trans_keys [] byte  = [] byte  { 32, 48, 57, 48, 57, 10, 32, 9, 13, 48, 57, 10, 32, 9, 13, 32, 48, 57, 32, 48, 57, 32, 0 }
var _cliparser_single_lengths [] int8  = [] int8  { 0, 1, 0, 2, 2, 1, 1, 1, 0, 0  }
var _cliparser_range_lengths [] int8  = [] int8  { 0, 1, 1, 2, 1, 1, 1, 0, 0, 0  }
var _cliparser_index_offsets [] int8  = [] int8  { 0, 0, 3, 5, 10, 14, 17, 20, 22, 0  }
var _cliparser_trans_cond_spaces [] int8  = [] int8  { -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 0  }
var _cliparser_trans_offsets [] int8  = [] int8  { 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 0  }
var _cliparser_trans_lengths [] int8  = [] int8  { 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0  }
var _cliparser_cond_keys [] int8  = [] int8  { 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0  }
var _cliparser_cond_targs [] int8  = [] int8  { 2, 5, 0, 3, 0, 8, 4, 4, 3, 0, 8, 4, 4, 0, 2, 6, 0, 2, 7, 0, 2, 0, 8, 0  }
var _cliparser_cond_actions [] int8  = [] int8  { 0, 1, 0, 3, 0, 5, 5, 5, 3, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 7, 0  }
var _cliparser_nfa_targs [] int8  = [] int8  { 0, 0  }
var _cliparser_nfa_offsets [] int8  = [] int8  { 0, 0, 0, 0, 0, 0, 0, 0, 0, 0  }
var _cliparser_nfa_push_actions [] int8  = [] int8  { 0, 0  }
var _cliparser_nfa_pop_trans [] int8  = [] int8  { 0, 0  }
var cliparser_start  int  = 1
var cliparser_first_final  int  = 8
var cliparser_error  int  = 0
var cliparser_en_main  int  = 1
type Cli struct {
	Status int
	Body []byte
}

func Cliparser(data []byte) (cli *Cli){
	cs, p, pe := 0, 0, len(data)
	cli = new(Cli)
	bodylength, bodypos := 0, 0
	
	{
		cs = int(cliparser_start);
	}
	
	{
		var  _klen int 
		var  _trans  uint   = 0
		var  _cond  uint   = 0
		var  _acts int
		var  _nacts uint 
		var  _keys int
		var  _ckeys int
		var  _cpc int 
		if p == pe  {
			goto _test_eof;
			
		}
		if cs == 0  {
			goto _out;
			
		}
		_resume :
		_keys = int(_cliparser_key_offsets[cs] );
		_trans = uint(_cliparser_index_offsets[cs]);
		_klen = int(_cliparser_single_lengths[cs]);
		if _klen > 0  {
			{
				var  _lower int
				var  _mid int
				var  _upper int
				_lower = _keys;
				_upper = _keys + _klen - 1;
				for {
					{
						if _upper < _lower  {
							break;
							
							
						}
						_mid = _lower + ((_upper-_lower) >> 1);
						switch {
							case ( data[p ]) < _cliparser_trans_keys[_mid ]:
							_upper = _mid - 1;
							
							case ( data[p ]) > _cliparser_trans_keys[_mid ]:
							_lower = _mid + 1;
							
							default:
							{
								_trans += uint((_mid - _keys));
								goto _match;
							}
							
						}
					}
					
				}
				_keys += _klen;
				_trans += uint(_klen);
			}
			
			
		}
		_klen = int(_cliparser_range_lengths[cs]);
		if _klen > 0  {
			{
				var  _lower int
				var  _mid int
				var  _upper int
				_lower = _keys;
				_upper = _keys + (_klen<<1) - 2;
				for {
					{
						if _upper < _lower  {
							break;
							
							
						}
						_mid = _lower + (((_upper-_lower) >> 1) & ^1);
						switch {
							case ( data[p ]) < _cliparser_trans_keys[_mid ]:
							_upper = _mid - 2;
							
							case ( data[p ]) > _cliparser_trans_keys[_mid + 1 ]:
							_lower = _mid + 2;
							
							default:
							{
								_trans += uint(((_mid - _keys)>>1));
								goto _match;
							}
							
						}
					}
					
				}
				_trans += uint(_klen);
			}
			
			
		}
		
		_match :
		_ckeys = int(_cliparser_trans_offsets[_trans] );
		_klen = int(_cliparser_trans_lengths[_trans]);
		_cond = uint(_cliparser_trans_offsets[_trans]);
		_cpc = 0;
		{
			var  _lower int
			var  _mid int
			var  _upper int
			_lower = _ckeys;
			_upper = _ckeys + _klen - 1;
			for {
				{
					if _upper < _lower  {
						break;
						
						
					}
					_mid = _lower + ((_upper-_lower) >> 1);
					switch {
						case _cpc < int(_cliparser_cond_keys[_mid ]):
						_upper = _mid - 1;
						
						case _cpc > int(_cliparser_cond_keys[_mid ]):
						_lower = _mid + 1;
						
						default:
						{
							_cond += uint((_mid - _ckeys));
							goto _match_cond;
						}
						
					}
				}
				
			}
			cs = 0;
			goto _again;
		}
		
		_match_cond :
		cs = int(_cliparser_cond_targs[_cond]);
		if _cliparser_cond_actions[_cond] == 0  {
			goto _again;
			
			
		}
		_acts = int(_cliparser_cond_actions[_cond] );
		_nacts = uint(_cliparser_actions[_acts ]);
		_acts += 1;
		for _nacts > 0  {
			{
				switch _cliparser_actions[_acts ] {
					case 0 :
					{cli.Status = cli.Status*10+(int((( data[p ])))-'0')}
					
					break;
					case 1 :
					{bodylength = bodylength*10+(int((( data[p ])))-'0')}
					
					break;
					case 2 :
					{cli.Body = make([]byte,bodylength)}
					
					break;
					case 3 :
					{if bodypos == bodylength {{p+= 1;
								goto _out; }}; cli.Body[bodypos]=(( data[p ])); bodypos++}
					
					break;
					
				}
				_nacts -= 1;
				_acts += 1;
			}
			
			
			
		}
		
		_again :
		if cs == 0  {
			goto _out;
			
		}
		p += 1;
		if p != pe  {
			goto _resume;
			
		}
		
		_test_eof :
		{}
		
		_out :
		{}
		
	}
	return cli
}

