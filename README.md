# Iknore File Server

#

example.png?width=100&height=100

example.png?width=100&height=100&cover=contain
example.png?width=100&height=100&cover=inside
example.png?width=100&height=100&cover=stretch

example.png?width=100&height=100&cover=contain

# 設定檔案

```yml
types:
    user_avatars:
        placeholder: "./placeholders/user.png"
        covers:
            - contain
        background_colors:
            - white
        sizes:
            - 200x
            - 400x (small)
            - 600x
        formats:
            - jpg
            - png
            - gif
        format_options:
            jpg:
                related_sizes: [200x, 400x, 600x]
                compression: COMPRESSION_B44A
                quality: 30
    circles:
        covers:
            - "*"
        background_colors:
            - "*"
        sizes:
            - "*"
placeholder: "./placeholders/default.png"
format_options:
    jpg:
        related_sizes: [200x, 400x, 600x]
        compression: COMPRESSION_B44A
        quality: 30
```
