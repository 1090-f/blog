package database

import "testing"

func TestDefaultSampleArticlesArePublishable(t *testing.T) {
	articles := defaultSampleArticles()
	if len(articles) < 3 {
		t.Fatalf("expected at least 3 sample articles, got %d", len(articles))
	}

	titles := make(map[string]struct{}, len(articles))
	for _, article := range articles {
		if article.Title == "" || article.Summary == "" || article.Content == "" {
			t.Fatal("sample article title, summary, and content must not be empty")
		}
		if article.CategoryName == "" || len(article.TagNames) == 0 {
			t.Fatal("sample article must have a category and at least one tag")
		}
		if _, exists := titles[article.Title]; exists {
			t.Fatalf("duplicate sample article title %q", article.Title)
		}
		titles[article.Title] = struct{}{}
	}
}
