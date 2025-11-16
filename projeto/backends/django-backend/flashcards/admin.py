from django.contrib import admin
from flashcards.models import Collection, Flashcard, Share, GenerationLog, Payment, User

@admin.register(User)
class UserAdmin(admin.ModelAdmin):
    list_display = ['email', 'username', 'plan', 'created_at']
    list_filter = ['plan', 'created_at']
    search_fields = ['email', 'username']

@admin.register(Collection)
class CollectionAdmin(admin.ModelAdmin):
    list_display = ['name', 'user', 'is_public', 'created_at']
    list_filter = ['is_public', 'created_at']
    search_fields = ['name', 'user__email']

@admin.register(Flashcard)
class FlashcardAdmin(admin.ModelAdmin):
    list_display = ['front', 'collection', 'created_by_ia', 'created_at']
    list_filter = ['created_by_ia', 'created_at']
    search_fields = ['front', 'back', 'collection__name']

@admin.register(Share)
class ShareAdmin(admin.ModelAdmin):
    list_display = ['collection', 'shared_with', 'permissions', 'created_at']
    list_filter = ['permissions', 'created_at']

@admin.register(GenerationLog)
class GenerationLogAdmin(admin.ModelAdmin):
    list_display = ['user', 'date', 'count']
    list_filter = ['date']

@admin.register(Payment)
class PaymentAdmin(admin.ModelAdmin):
    list_display = ['user', 'status', 'created_at']
    list_filter = ['status', 'created_at']