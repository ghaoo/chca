{{ define "content" }}

<div class="post-page">
    <div class="post animated fadeInDown">
        <div class="post-title">
            <h3>{{ .article.Title }}
            </h3>
        </div>
        <div class="post-content" id="content">

            {{ .article.Content | unescaped }}

            <nav class="article-nav" id="state">
                <span class="label label-important">PERMANENT LINK:</span>
                <a href="http://chca.me{{.article.Url}}">http://chca.me{{.article.Url}}</a>
            </nav>

        </div>

        <div class="post-footer">
            <div class="meta">
                <div class="info">
                    <i class="fa fa-calendar"></i>
                    <span class="date">{{ .article.CreatedAt | format }}</span>
                    &nbsp;&nbsp;
                    <i class="fa fa-folder-open"></i>
                    {{ range .article.Category }}
                    <a href="/category/{{ . }}">{{ . }}</a>&nbsp;
                    {{ end }}
                    &nbsp;&nbsp;
                    <i class="fa fa-tags"></i>
                    {{ range .article.Tags }}
                    <a href="/tag/{{ . }}">{{ . }}</a>&nbsp;
                    {{ end }}
                </div>
            </div>
        </div>
    </div>

    <!--评论-->

</div>

{{ end }}