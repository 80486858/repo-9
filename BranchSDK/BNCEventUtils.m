//
//  BNCEventUtils.m
//  BranchSDK
//
//  Created by Nipun Singh on 1/31/23.
//  Copyright © 2023 Branch, Inc. All rights reserved.
//

#import "BNCEventUtils.h"

@interface BNCEventUtils()
@property (nonatomic, strong, readwrite) NSMutableSet *events;
@end

@implementation BNCEventUtils

+ (instancetype)shared {
    static BNCEventUtils *set;
    static dispatch_once_t onceToken;
    dispatch_once(&onceToken, ^{
        set = [BNCEventUtils new];
    });
    return set;
}

- (instancetype)init {
    self = [super init];
    if (self) {
        self.events = [NSMutableSet alloc];
    }
    return self;
}

- (void)storeEvent:(BranchEvent *)event withCompletion:(void (^_Nullable)(BOOL success, NSError * _Nullable error))completion {
    
    [self.events addObject:event];
}

- (void)removeEvent:(BranchEvent *)event withCompletion:(void (^_Nullable)(BOOL success, NSError * _Nullable error))completion {
    [self.events removeObject:event];
}

@end
