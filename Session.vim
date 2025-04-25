let SessionLoad = 1
let s:so_save = &g:so | let s:siso_save = &g:siso | setg so=0 siso=0 | setl so=-1 siso=-1
let v:this_session=expand("<sfile>:p")
silent only
silent tabonly
cd ~/projects/crypto-trading-bot
if expand('%') == '' && !&modified && line('$') <= 1 && getline(1) == ''
  let s:wipebuf = bufnr('%')
endif
let s:shortmess_save = &shortmess
if &shortmess =~ 'A'
  set shortmess=aoOA
else
  set shortmess=aoO
endif
badd +396 internal/services/marketdata_service.go
badd +56 internal/models/marketdata.go
badd +1 data/waves.json
badd +0 ~/.config/nvim/init.vim
argglobal
%argdel
$argadd ~/.config/nvim/init.vim
edit internal/services/marketdata_service.go
argglobal
balt data/waves.json
setlocal fdm=indent
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
30
normal! zo
191
normal! zo
224
normal! zo
225
normal! zo
235
normal! zo
286
normal! zo
290
normal! zo
299
normal! zo
300
normal! zo
323
normal! zo
324
normal! zo
337
normal! zo
338
normal! zo
343
normal! zo
361
normal! zo
362
normal! zo
376
normal! zo
379
normal! zo
380
normal! zo
387
normal! zo
388
normal! zo
395
normal! zo
397
normal! zo
401
normal! zo
419
normal! zo
421
normal! zo
424
normal! zo
429
normal! zo
435
normal! zo
447
normal! zo
454
normal! zo
472
normal! zo
475
normal! zo
486
normal! zo
491
normal! zo
508
normal! zo
let s:l = 396 - ((11 * winheight(0) + 14) / 28)
if s:l < 1 | let s:l = 1 | endif
keepjumps exe s:l
normal! zt
keepjumps 396
normal! 033|
tabnext 1
if exists('s:wipebuf') && len(win_findbuf(s:wipebuf)) == 0 && getbufvar(s:wipebuf, '&buftype') isnot# 'terminal'
  silent exe 'bwipe ' . s:wipebuf
endif
unlet! s:wipebuf
set winheight=1 winwidth=20
let &shortmess = s:shortmess_save
let s:sx = expand("<sfile>:p:r")."x.vim"
if filereadable(s:sx)
  exe "source " . fnameescape(s:sx)
endif
let &g:so = s:so_save | let &g:siso = s:siso_save
set hlsearch
doautoall SessionLoadPost
unlet SessionLoad
" vim: set ft=vim :
