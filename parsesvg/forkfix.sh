#!/bin/zsh

#https://stackoverflow.com/questions/11392478/how-to-replace-a-string-in-multiple-files-in-linux-command-line

grep -rli 'timdrysdale' * | xargs -i@ sed -i 's/timdrysdale/timdrysdale/g' @
